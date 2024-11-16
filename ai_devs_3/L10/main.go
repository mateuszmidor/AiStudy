package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
	"github.com/pkg/errors"
)

const filesZipURL = "https://centrala.ag3nts.org/dane/pliki_z_fabryki.zip"

type FileCategories struct {
	People   []string `json:"people"`
	Hardware []string `json:"hardware"`
}

func downloadZipIfDoesntExistYet(url, destination string) error {
	// Check if the destination file already exists
	if _, err := os.Stat(destination); err == nil {
		// File already exists, no need to download
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than "file does not exist" occurred
		return errors.Wrap(err, "failed to check if destination file exists")
	}

	// Create the downloads directory if it doesn't exist
	downloadsDir := filepath.Dir(destination)
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadsDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed to create downloads directory")
		}
	}

	// Download the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	// Create the destination file to store the downloaded ZIP
	destFile, err := os.Create(destination)
	if err != nil {
		return errors.Wrap(err, "failed to create destination file")
	}
	defer destFile.Close()

	// Write the response body to the destination file
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to write to destination file")
	}

	return nil
}

// readOrTranscribeAudio returns key-value pairs: {filename: transcription} for all .mp3 files found under sourceDir
func readOrTranscribeAudio(sourceDir string) (map[string]string, error) {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	transcriptions := make(map[string]string)

	for _, file := range files {
		// skip non-audio files
		if strings.ToLower(filepath.Ext(file.Name())) != ".mp3" {
			continue
		}

		// just read the transcription if already exists
		txtFileName := filepath.Join(sourceDir, file.Name()+".txt")
		if _, err := os.Stat(txtFileName); err == nil {
			content, err := os.ReadFile(txtFileName)
			if err != nil {
				log.Printf("failed to read transcription file %s: %+v\n", txtFileName, err)
			} else {
				transcriptions[file.Name()] = string(content)
			}
			continue
		}

		// do transcribe
		audioFileName := filepath.Join(sourceDir, file.Name())
		log.Println("transcribing", audioFileName)
		text, err := openai.SpeachToText(audioFileName)
		if err != nil {
			log.Printf("failed: %+v\n", err)
		}
		transcriptions[file.Name()] = text

		// save transcription for next-time use
		err = os.WriteFile(txtFileName, []byte(text), os.ModePerm)
		if err != nil {
			log.Printf("failed to save transcription to file %s: %+v\n", txtFileName, err)
		}
	}

	return transcriptions, nil
}

// readOrTranscribeImages returns key-value pairs: {filename: transcription} for all .png files found under sourceDir
func readOrTranscribeImages(sourceDir string) (map[string]string, error) {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	transcriptions := make(map[string]string)

	for _, file := range files {
		// skip non-image files
		if strings.ToLower(filepath.Ext(file.Name())) != ".png" {
			continue
		}

		// just read the transcription if already exists
		txtFileName := filepath.Join(sourceDir, file.Name()+".txt")
		if _, err := os.Stat(txtFileName); err == nil {
			content, err := os.ReadFile(txtFileName)
			if err != nil {
				log.Printf("failed to read transcription file %s: %+v\n", txtFileName, err)
			} else {
				transcriptions[file.Name()] = string(content)
			}
			continue
		}

		// do transcribe
		imageFileName := filepath.Join(sourceDir, file.Name())
		log.Println("transcribing", imageFileName)
		system := "Read and return all the text from the attached image"
		user := "Return only the text from the image, WHITOUT any additional comments"
		text, err := openai.CompletionCheap(user, system, []string{openai.ImageFromFile(imageFileName)})
		if err != nil {
			log.Printf("failed: %+v\n", err)
		}
		transcriptions[file.Name()] = text

		// save transcription for next-time use
		err = os.WriteFile(txtFileName, []byte(text), os.ModePerm)
		if err != nil {
			log.Printf("failed to save transcription to file %s: %+v\n", txtFileName, err)
		}
	}

	return transcriptions, nil
}

// buid single text document from multiple .txt files.
// Example output:
// [filename1]
// file content 1

// [filename2]
// file content 2

// [filename3]
// file content 3
func buildCollectiveDocumentFromTxtFiles(srcDir string) string {

	var allDocuments string

	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("failed to read directory: %+v", err)
	}
	for _, file := range files {

		if strings.ToLower(filepath.Ext(file.Name())) != ".txt" {
			continue
		}

		originalFileName := removeExtraExt(file.Name())

		filePath := filepath.Join(srcDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", filePath, err)
			continue
		}

		allDocuments += fmt.Sprintf("[%s]\n%s\n\n", originalFileName, string(content))
	}
	return allDocuments
}

func categorizeFilesInCollectiveDocument(documentWithMultipleFiles string) FileCategories {
	const system = `
	You will be given a text document composed of a series FileName-FileContent blocks. Example input document:
	<example_input>
	[filename1]
	file content 1

	[filename2]
	file content 2

	[filename3]
	file content 3
	</example_input>
	Your role is to categorize content of each file to be either:
	1. people-related (only people which were caught or seen) 
	2. hardware-related (only hardware which was physically repaired, not the hardware that had software updated)
	3. unrelated to people or hardware
	Return the names of files that are either people-related or hardware-related in form of JSON, ignore the unrelated files. Example JSON result:
	<example_output>
	{
		"people": ["plik1.txt", "plik2.mp3", "plikN.png"],
		"hardware": ["plik4.txt", "plik5.png", "plik6.mp3"],
	}
	</example_output>
	`
	user := "The document: \n" + documentWithMultipleFiles
	resultString, err := openai.CompletionExpert(user, system, nil, "gpt-4o-mini", "json_object", 1000, 0.0)
	if err != nil {
		log.Fatalf("openai returend error: %+v", err)
	}

	fmt.Println(resultString.Error)
	result := FileCategories{}
	err = json.Unmarshal([]byte(resultString.Choices[0].Message.Content), &result)
	if err != nil {
		log.Fatalf("Error deserializing JSON: %+v", err)
	}
	return result
}

func removeExtraExt(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	if strings.Contains(base, ".") {
		return base
	}
	return filename
}

func splitIdQuestion(questions string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(questions, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result

}

func fixHyperlinks(document string) string {
	// return strings.ReplaceAll(document, "(i/", "(https://centrala.ag3nts.org/dane/i/")
	return document
}

func extractLinks(document string) []string {
	var links []string
	re := regexp.MustCompile(`\((i/[^)]+)\)`)
	matches := re.FindAllStringSubmatch(document, -1)
	for _, match := range matches {
		if len(match) > 1 {
			links = append(links, match[1])
		}
	}
	return links
}

func downloadFromHyperlinks(links []string) error {
	for _, link := range links {
		// Append the prefix to the link
		fullURL := "https://centrala.ag3nts.org/dane/" + link

		// Create the downloads directory if it doesn't exist
		downloadsDir := "downloads"
		if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
			err = os.Mkdir(downloadsDir, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "failed to create downloads directory")
			}
		}

		// Determine the destination file path
		fileName := filepath.Base(link)
		destination := filepath.Join(downloadsDir, fileName)

		// Download the file
		resp, err := http.Get(fullURL)
		if err != nil {
			return errors.Wrapf(err, "failed to download file from %s", fullURL)
		}
		defer resp.Body.Close()

		// Create the destination file
		destFile, err := os.Create(destination)
		if err != nil {
			return errors.Wrapf(err, "failed to create destination file %s", destination)
		}
		defer destFile.Close()

		// Write the response body to the destination file
		_, err = io.Copy(destFile, resp.Body)
		if err != nil {
			return errors.Wrapf(err, "failed to write to destination file %s", destination)
		}
	}
	return nil
}

func main() {
	// fetch the questions
	apiKey := api.ApiKey()
	questionsURL := "https://centrala.ag3nts.org/data/" + apiKey + "/arxiv.txt"
	questions, err := api.FetchData(questionsURL)
	if err != nil {
		log.Fatalf("Error fetching questions: %+v", err)
	}
	idQuestionMap := splitIdQuestion(questions)
	_ = idQuestionMap

	// fetch the article
	articleURL := "https://centrala.ag3nts.org/dane/arxiv-draft.html"
	articleHTML, err := api.FetchData(articleURL)
	if err != nil {
		log.Fatalf("Error fetching questions: %+v", err)
	}

	// transform article to MarkDown for LLM to understand better
	articleMD, err := htmltomarkdown.ConvertString(articleHTML)
	if err != nil {
		log.Fatal(err)
	}
	fixedArticleMD := fixHyperlinks(articleMD)
	err = os.WriteFile("article.md", []byte(fixedArticleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to file article.md: %+v", err)
	}

	// extract all multimedia links
	multimediaHyperlinks := extractLinks(articleMD)
	downloadFromHyperlinks(multimediaHyperlinks)
}
