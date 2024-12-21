package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

// rebuild the dialogs one by one from shortest one to longest and return the rebuilt dialogs
func rebuildDialogs(dialogs Root) map[string]string {
	result := map[string]string{} // filename -> reconstructed dialog

	// collect fragmented dialogs into a list and sort by length ascending to improve chances for success - easier to reconstruct short dialogs first
	fragmentedDialogs := []Rozmowa{dialogs.Rozmowa1, dialogs.Rozmowa2, dialogs.Rozmowa3, dialogs.Rozmowa4, dialogs.Rozmowa5}
	slices.SortFunc(fragmentedDialogs, func(a, b Rozmowa) int {
		return a.Length - b.Length
	})

	// rebuild dialogs one by one, store rebuilt dialogs in files for the next time use
	dialogPieces := dialogs.Reszta
	for i, d := range fragmentedDialogs {
		filename := fmt.Sprintf("%d.txt", i)

		// if dialog is already rebuilt - read it in and move to the next one
		if dialogBytes, err := os.ReadFile(filename); err == nil {
			dialog := string(dialogBytes)
			log.Printf("found %s; reading in and moving on to the next one", filename)
			result[filename] = dialog
			dialogPieces = deleteUsedPieces(dialogPieces, dialog)
			continue
		}

		// try to rebuild the dialog. repeat until success :)
		var dialog string
		originalPieces := dialogPieces
		for {
			dialog, dialogPieces = rebuildDialog(d, dialogPieces)
			if dialog != "" {
				break // success
			}
			dialogPieces = originalPieces // restore pieces for another attempt
		}
		result[filename] = dialog

		// save the rebuilt dialog
		if err := os.WriteFile(filename, []byte(dialog), 0644); err != nil {
			log.Printf("failed to save %s: %v", filename, err)
		} else {
			log.Printf("successfully saved %s", filename)
		}
	}

	return result
}

func rebuildDialog(fragmentedDialog Rozmowa, pieces []string) (string, []string) {
	const system = `Od teraz twoja rola polega na odbudowaniu kompletnej rozmowy dwóch osób z wymieszanych fragmentów (kwestii).
Znany jest koniec rozmowy, podany w sekcji <Koniec Rozmowy>. Cała rozmowa musi zmierzać do kwesti w sekcji <Koniec Rozmowy>.
Znana jest całkowita długość rozmowy (liczba wypowiedzianych kwestii w trakcie całej rozmowy) w sekcji <Długość Rozmowy>.
Dopasuj kolejną najlepiej pasującą kwestię ze zbioru <Wymieszane Kwestie> jako kontynuację rozmowy <Dotychczasowa Rozmowa>. 
Wynikowa rozmowa ma być logiczna i zmierzać do kwestii <Koniec Rozmowy>, nie przekraczając <Długość Rozmowy>.
Pamiętaj, ze potrzebne są tylko niektóre kwestie ze zbioru <Wymieszane Kwestie>, nie będziesz musiał uzyc wszystkich.
Zwróć jedynie samą dopasowaną kwestię, bez komentarzy ani formatowania.
Wazne: Nie wolno Ci uzywac kwestii juz wypowidzianych w sekcji <Dotychczasowa Rozmowa>. Uzywaj tylko kwestii dostepnych w sekcji <Wymieszane Kwestie>.
Przykład:
<Wymieszane Kwestie>
"- Wiedziałeś ze Barbara urodziła się w kwietniu?"
"- Dziś gwiazdy świecą mocniej niz zazwyczaj"
"- Biegałem, bo lubię zacząć dzień od sportu"
<Wymieszane Kwestie/>
<Dotychczasowa Rozmowa>
"- Co robiłeś wczoraj rano?"
<Dotychczasowa Rozmowa/>
<Koniec Rozmowy>
"- Aha, dawka ruchu o poranku to dobry pomysł"
<Koniec Rozmowy/>
<Długość Rozmowy>
3
<Długość Rozmowy/>

<Dopasowana Kwestia>
"- Biegałem, bo lubię zacząć dzień od sportu"
<Dopasowana Kwestia/>
`
	const userFormat = `<Wymieszane Kwestie>
%s
<Wymieszane Kwestie/>
<Koniec Rozmowy>
%s
<Koniec Rozmowy/>
<Długość Rozmowy>
%d
<Długość Rozmowy/>
<Dotychczasowa Rozmowa>
%s
<Dotychczasowa Rozmowa/>"`

	conversation := fragmentedDialog.Start
	for i := 1; i <= fragmentedDialog.Length; i++ { // +2 for begin and end
		fmt.Printf("\nnum pieces: %d, step %d/%d\n", len(pieces), i, fragmentedDialog.Length)
		fmt.Print("Press 'Enter' to continue...")
		fmt.Scanln()
		user := fmt.Sprintf(userFormat, strings.Join(pieces, "\n"), fragmentedDialog.End, fragmentedDialog.Length, conversation)
		fmt.Println(user)
		rsp, err := openai.CompletionStrong(user, system, nil)
		if err != nil {
			log.Fatal(err)
		}
		text := rsp
		// fmt.Println("### rsp:", text)
		if text == fragmentedDialog.End {
			fmt.Println(text)
			log.Print("FINISH")
			return conversation, pieces
		}
		// Remove the selected text from pieces using slices package
		pieces = deleteUsedPieces(pieces, text)
		conversation += "\n" + text
	}
	fmt.Println("FAIL")
	return "", pieces
}

func deleteUsedPieces(pieces []string, text string) []string {
	return slices.DeleteFunc(pieces, func(piece string) bool {
		return strings.Contains(text, piece)
	})
}
