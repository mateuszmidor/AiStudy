package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"

	"golang.org/x/text/unicode/norm"
)

type Root struct {
	Rozmowa1 Rozmowa1 `json:"rozmowa1"`
	Rozmowa2 Rozmowa2 `json:"rozmowa2"`
	Rozmowa3 Rozmowa3 `json:"rozmowa3"`
	Rozmowa4 Rozmowa4 `json:"rozmowa4"`
	Rozmowa5 Rozmowa5 `json:"rozmowa5"`
	Reszta   []string `json:"reszta"`
}
type DialogInfo struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Length int    `json:"length"`
}

type Rozmowa1 struct {
	DialogInfo
}
type Rozmowa2 struct {
	DialogInfo
}
type Rozmowa3 struct {
	DialogInfo
}
type Rozmowa4 struct {
	DialogInfo
}
type Rozmowa5 struct {
	DialogInfo
}

func main() {
	// get questions
	questionsURL := "https://centrala.ag3nts.org/data/" + api.ApiKey() + "/phone_questions.json"
	questions := map[string]string{}
	api.FillDataFromJSONURL(questionsURL, &questions)
	// fmt.Println(questions)

	// get dialogs
	dialogsURL := "https://centrala.ag3nts.org/data/" + api.ApiKey() + "/phone.json"
	var phone = Root{}
	api.FillDataFromJSONURL(dialogsURL, &phone)

	// Fix some weird and problematic whitespace in the data
	fixText(&phone)

	// Dump dialogs for troubleshooting
	phoneData, err := json.MarshalIndent(phone, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("phone.json", phoneData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// get additional facts
	// facts := composeDocumentFromFragments("facts/")
	// fmt.Println(facts)

	rebuildDialogs(phone)
}

// input text contains '\u00A0' that complicates text processing; let's remove it
func fixText(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	dataStr := string(data)
	newDataStr := strings.ReplaceAll(dataStr, "\u00A0", " ")
	newData := []byte(newDataStr)
	err = json.Unmarshal(newData, &v)
	if err != nil {
		log.Fatal(err)
	}
}

func rebuildDialogs(dialogs Root) map[string]string {

	result := map[string]string{}

	rebuildDialog(dialogs.Rozmowa3.DialogInfo, dialogs.Reszta)
	return result
}

func rebuildDialog(info DialogInfo, pieces []string) {
	const system = `
Dopasuj najlepiej pasującą kwestię spośród <Dostępne Kwestie> jako kontynuację <Dotychczasowa Rozmowa>. Zwróć jedynie samą dopasowaną kwestię, bez komentarzy ani formatowania.
Nie wolno Ci uzywac kwestii juz wypowidzianych w sekcji <Dotychczasowa Rozmowa>. Uzywaj tylko kwestii dostepnych w sekcji <Dostępne Kwestie>.
Przykład:
<Dostępne Kwestie>
"- Wiedziałeś ze Barbara urodziła się w kwietniu?"
"- Dziś gwiazdy świecą mocniej niz zazwyczaj"
"- Biegałem, bo lubię zacząć dzień od sportu"
<Dostępne Kwestie/>
<Dotychczasowa Rozmowa>
"- Co robiłeś wczoraj rano?"
<Dotychczasowa Rozmowa/>
<Dopasowanie>
"- Biegałem, bo lubię zacząć dzień od sportu"
<Dopasowanie>
`

	pieces = append(pieces, info.End)
	for i := range pieces {
		pieces[i] = norm.NFC.String(pieces[i])
	}

	lines := info.Start
	for i := 1; i <= info.Length+2; i++ { // +2 for begin and end
		fmt.Printf("\nnum pieces: %d, step %d/%d\n", len(pieces), i, info.Length)
		fmt.Print("Press 'Enter' to continue...")
		fmt.Scanln()
		user := "<Dostępne Kwestie>\n" + strings.Join(pieces, "\n") + "\n<Dostępne Kwestie/>\n<Dotychczasowa Rozmowa>\n" + lines + "\n<Dotychczasowa Rozmowa/>"

		fmt.Println(user)
		rsp, err := openai.CompletionStrong(user, system, nil)
		if err != nil {
			log.Fatal(err)
		}
		text := rsp
		// fmt.Println("### rsp:", text)
		if text == info.End {
			fmt.Println(text)
			log.Print("FINISH")
			return
		}
		// Remove the selected text from pieces using slices package
		pieces = slices.DeleteFunc(pieces, func(piece string) bool {
			return piece == text
		})
		lines += "\n" + text
	}
	fmt.Println("FAIL")
}

func composeDocumentFromFragments(sourceDir string) string {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Printf("failed to read directory: %+v\n", err)
		return ""
	}

	var fragments []string

	// Collect all .txt files with the expected naming pattern
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			fragments = append(fragments, file.Name())
		}
	}

	// Sort the files in numeric order based on the number in the filename
	sort.Slice(fragments, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fragments[i], "f-"), ".txt"))
		numJ, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fragments[j], "f-"), ".txt"))
		return numI < numJ
	})

	var documentBuilder strings.Builder

	// Read and append the content of each file to the document
	for _, fragment := range fragments {
		content, err := os.ReadFile(filepath.Join(sourceDir, fragment))
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", fragment, err)
			continue
		}
		documentBuilder.Write(content)
		documentBuilder.WriteString("\n") // Add a newline between fragments
	}

	doc := documentBuilder.String()
	err = os.WriteFile("document.txt", []byte(doc), os.ModePerm)
	if err != nil {
		log.Printf("failed to save document to file document.txt: %+v\n", err)
	}
	return doc
}

// - tak Zygfryd, słyszę Cię teraz dobrze. Przepraszam, gdy poprzednio dzwoniłeś, byłem w fabryce. Wiesz, w sektorze D, gdzie się produkuje broń i tutaj mają jakąś izolację na ścianach dodatkową. Telefon gubi zasięg. Masz jakieś nowe zadanie dla mnie?
// - tak Zygfryd, słyszę Cię teraz dobrze. Przepraszam, gdy poprzednio dzwoniłeś, byłem w fabryce. Wiesz, w sektorze D, gdzie się produkuje broń i tutaj mają jakąś izolację na ścianach dodatkową. Telefon gubi zasięg. Masz jakieś nowe zadanie dla mnie?
