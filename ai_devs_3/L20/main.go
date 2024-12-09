package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath" // Add this line
	"sort"
	"strconv"

	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
	"github.com/pkg/errors"
)

const system = `
# Kontekst:
Otrzymasz tekst źródłowy na podstawie którego odpowiesz na pytanie uzytkownika.
Tekst źródłowy ma postać osobistego notatnika Rafała - jest on napisany w sposób nieuporządkownay, wyrywkowo i niedbale, Rafał opisuje w nim swoje przemyślenia
i wspomina wydarzenia z ostatnich lat.
# Zadanie:
Odpowiedz JEDNYM SŁOWEM na pytanie uzytkownika, odpowiedź ma być w formacie JSON z polami "_przemyślenia", "błędne_odpowiedzi" oraz "odpowiedź", gdzie pole "_przemyślenia" występuje jako pierwsze.
Zacznij od zgromadzenia wszystkich informacji, które mogą być pomocne by odpowiedzieć na pytanie, i umieść je w polu "_przemyślenia", nie pomijaj niczego!
Następnie na podstawie zgromadzonych informacji wydedukuj najbardziej prawdopodobną odpowiedź i umieść ją w polu "answer".
Jeśli uzytkownik powie, ze odpowiedź jest błędna, zapamiętaj ją w polu "wrong_answers", tak zeby nie podać błędnej odpowiedzi przy kolejnej próbie.
# Przykład:
pytanie: w którym roku Rafał spotkał Andrzeja?
odpowiedź: 
{
"_przemyślenia":"w notatniku Rafał wspomina, ze odwiedził Walencję w 2015 roku i ze było to 6 lat wcześniej, niz poznał Andrzeja",
"błędne_odpowiedzi":["2024","2023","2022"],
"odpowiedź":"2021"
}
Pamiętaj - odpowiadaj jednym słowem.
# Notatnik Rafała:
`

type Answer struct {
	Thinking     string   `json:"_przemyślenia"`
	WrongAnswers []string `json:"błędne_odpowiedzi"`
	Answer       string   `json:"odpowiedź"`
}

func main() {
	_, err := readOrDescribeImages("downloads/")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	doc := composeDocumentFromFragments("downloads/")

	url := "https://centrala.ag3nts.org/data/" + api.ApiKey() + "/notes.json"
	questionsJSON, err := api.FetchData(url)
	if err != nil {
		log.Fatal(err)
	}
	questions := map[string]string{}
	err = json.Unmarshal([]byte(questionsJSON), &questions)
	if err != nil {
		log.Fatal("Error deserializing questions JSON:", err)
	}

	// fill answers with keys so API responds with wrong answer ID so we know we must find another answer
	answers := make(map[string]string, len(questions))
	for id := range questions {
		answers[id] = ""
	}

	// Assign sorted keys from questions - so we answer questions 1,2,3,4,5 and get feedback from API
	ids := make([]string, 0, len(questions))
	for id := range questions {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		// get question
		user := questions[id]

		// try answer the question
		chat, err := openai.NewChatWithMemory(system+doc, "gpt-4o-mini", 1000, true)
		if err != nil {
			log.Fatal(err)
		}
		rsp, err := chat.User(user, nil, nil, "json_object", 0)
		if err != nil {
			log.Fatal(err)
		}

		// repeat until we get correct answer for this question
		for {
			fmt.Println("Press 'Enter' to continue...")
			fmt.Scanln()

			// verify answer
			var answer Answer
			err = json.Unmarshal([]byte(rsp.Choices[0].Message.Content), &answer)
			if err != nil {
				log.Fatal("Error deserializing response content:", err)
			}
			answers[id] = answer.Answer
			result, err := api.VerifyTaskAnswer("notes", answers, api.VerificationURL)
			if err != nil {
				if strings.Contains(err.Error(), id+" is incorrect") {
					// wrong answer, try again
					rsp, err = chat.User("Błędna odpowiedź. Zastanów się chwilę i podaj poprawną odpowiedź. NIE PODAWAJ TEJ SAMEJ BŁĘDNEJ ODPOWIEDZI PONOWNIE!", nil, nil, "json_object", 0) // feedback error to the chat to try again
					if err != nil {
						log.Fatal(err)
					}
				} else {
					// good answer! move on to he next question
					break
				}
			} else {
				// victory
				fmt.Println(result)
				return
			}
		}
	}
	// answers["01"] = "2019"
	// answers["02"] = "Adam"
	// answers["03"] = "Jaskinia"
	// answers["04"] = "2024-11-12"
	// answers["05"] = "Lubawa"
}

// readOrDescribeImages returns key-value pairs: {filename: description} for all .png files found under sourceDir
func readOrDescribeImages(sourceDir string) (map[string]string, error) {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	descriptions := make(map[string]string)

	for _, file := range files {
		// skip non-image files
		if strings.ToLower(filepath.Ext(file.Name())) != ".png" {
			continue
		}

		// just read the description if already exists
		txtFileName := filepath.Join(sourceDir, file.Name()+".txt")
		if _, err := os.Stat(txtFileName); err == nil {
			content, err := os.ReadFile(txtFileName)
			if err != nil {
				log.Printf("failed to read description file %s: %+v\n", txtFileName, err)
			} else {
				descriptions[file.Name()] = string(content)
			}
			continue
		}

		// do describe
		imageFileName := filepath.Join(sourceDir, file.Name())
		log.Println("describing", imageFileName)
		system := "Jesteś ekspertem OCR - potrafisz czytać tekst z obrazów"
		user := "Przepisz słowo w słowo tekst widoczny na obrazie. Jeśli są to fragmenty tekstu - tez je przepisz. Jeśli fragmenty tekstu są nieczytelne - domyśl się i uzupełnij je tak zeby pasowały do całości tekstu"
		text, err := openai.CompletionCheap(user, system, []string{openai.ImageFromFile(imageFileName)})
		if err != nil {
			log.Printf("failed: %+v\n", err)
		}
		descriptions[file.Name()] = text

		// save transcription for next-time use
		err = os.WriteFile(txtFileName, []byte(text), os.ModePerm)
		if err != nil {
			log.Printf("failed to save transcription to file %s: %+v\n", txtFileName, err)
		}
	}

	return descriptions, nil
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
		if strings.HasPrefix(file.Name(), "pdf-") && strings.HasSuffix(file.Name(), ".png.txt") {
			fragments = append(fragments, file.Name())
		}
	}

	// Sort the files in numeric order based on the number in the filename
	sort.Slice(fragments, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fragments[i], "pdf-"), ".png.txt"))
		numJ, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fragments[j], "pdf-"), ".png.txt"))
		return numI < numJ
	})

	var documentBuilder strings.Builder

	// Read and append the content of each file to the document
	for i, fragment := range fragments {
		content, err := os.ReadFile(filepath.Join(sourceDir, fragment))
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", fragment, err)
			continue
		}
		contentWithoutEmptyLinest := strings.ReplaceAll(string(content), "\n\n", "\n")
		documentBuilder.WriteString(fmt.Sprintf("Notatnik - strona %d\n", i+1))
		documentBuilder.WriteString(contentWithoutEmptyLinest)
		documentBuilder.WriteString("\n") // Add a newline between fragments
		documentBuilder.WriteString("\n") // Add a newline between fragments
	}

	doc := documentBuilder.String()
	err = os.WriteFile("document.txt", []byte(doc), os.ModePerm)
	if err != nil {
		log.Printf("failed to save document to file document.txt: %+v\n", err)
	}
	return doc
}
