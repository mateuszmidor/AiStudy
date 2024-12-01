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
# Zadanie
w oparciu o dostarczony Notatnik Rafała, odpowiedz na pytanie. Notatnik jest pisany zdawkowo i niedbale, odpowiedzi na pytanie trzeba starannie poszukać.
Odpowiadaj maksymalnie zwięźle i krótko, to wazne! 

# Przykład
	pytanie - w którym roku toczy wydarzenie miało miejsce?
	odpowiedź - 2123

# Notatnik Rafała
`

func main() {
	_, err := readOrDescribeImages("downloads/")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// fmt.Println(filenames)
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

	answers := map[string]string{}
	for id, question := range questions {
		user := question
		rsp, err := openai.CompletionCheap(user, system+doc, nil)
		if err != nil {
			log.Fatal(err)
		}
		answers[id] = rsp
	}

	// answers["01"] = "2019"
	// answers["02"] = "Adam"
	// answers["03"] = "Jaskinia"
	// answers["04"] = "2024-11-12"
	// answers["05"] = "Lubawa"

	for id, answer := range answers {
		fmt.Printf("%s: %s - %s\n", id, questions[id], answer)
	}
	result, err := api.VerifyTaskAnswer("notes", answers, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
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
