package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
)

type Root struct {
	Rozmowa1 Rozmowa  `json:"rozmowa1"`
	Rozmowa2 Rozmowa  `json:"rozmowa2"`
	Rozmowa3 Rozmowa  `json:"rozmowa3"`
	Rozmowa4 Rozmowa  `json:"rozmowa4"`
	Rozmowa5 Rozmowa  `json:"rozmowa5"`
	Reszta   []string `json:"reszta"`
}

type Rozmowa struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Length int    `json:"length"` // e.g. 3 means, Person1, Person2, Person3 - it includes Start and End
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
	fixJSON(&phone)

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
func fixJSON(v any) {
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
