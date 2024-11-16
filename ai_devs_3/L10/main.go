package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
	"github.com/pkg/errors"
)

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
		system := "You MUST use Polish language and Polish ONLY"
		user := "Describe what thing is on the picture, in case of a place, tell me the exact name of the place"
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

// input format: "01=gdzie zaczyna się wisła?""
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

// Examples:
// ![](i/rynek.png) -> i/rynek.png
// [rafal\_dyktafon.mp3](i/rafal_dyktafon.mp3) -> i/rafal_dyktafon.mp3
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

func downloadFromHyperlinks(filePaths []string) error {
	downloadsDir := "i" // article.md multimedia live in such folder, so converted article.md will work stright away
	for _, filePath := range filePaths {
		// Append the prefix to the link
		fileURL := "https://centrala.ag3nts.org/dane/" + filePath

		// Create the downloads directory if it doesn't exist
		if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
			err = os.Mkdir(downloadsDir, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "failed to create downloads directory")
			}
		}

		// Determine the destination file path
		fileName := filepath.Base(filePath)
		destination := filepath.Join(downloadsDir, fileName)

		// Download the file
		resp, err := http.Get(fileURL)
		if err != nil {
			return errors.Wrapf(err, "failed to download file from %s", fileURL)
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

func substiteDescriptionsForMultimedia(document string, descriptions map[string]string) string {
	// Regular expression to find multimedia links in the format ![](i/filename)
	re := regexp.MustCompile(`\!?\[.*\]\((.*)\)`)

	// Replace each multimedia link with its corresponding description
	result := re.ReplaceAllStringFunc(document, func(link string) string {
		// Extract the filename from the link
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			filename := filepath.Base(matches[1])
			// Check if there is a description for this filename
			if description, exists := descriptions[filename]; exists {
				// Return the description enclosed in triple backticks
				return "```\n" + description + "\n```"
			} else {
				fmt.Println("no description found for:", filename)
			}
		}
		// If no description is found, return the original link
		return link
	})

	return result
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

	// transform article to MarkDown for LLM to understand it better
	articleMD, err := htmltomarkdown.ConvertString(articleHTML)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("article.md", []byte(articleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to file article.md: %+v", err)
	}

	// extract all multimedia links and download files for transcrption
	multimediaHyperlinks := extractLinks(articleMD)
	downloadFromHyperlinks(multimediaHyperlinks)
	audioTranscriptions, err := readOrTranscribeAudio("i/")
	if err != nil {
		log.Fatal(err)
	}
	imageTranscriptions, err := readOrDescribeImages("i/")
	if err != nil {
		log.Fatal(err)
	}
	multimediaAsText := make(map[string]string)
	for k, v := range audioTranscriptions {
		multimediaAsText[k] = v
	}
	for k, v := range imageTranscriptions {
		multimediaAsText[k] = v
	}

	// substitute hyperlinks with media transcriptions
	substitutedArticleMD := substiteDescriptionsForMultimedia(articleMD, multimediaAsText)
	err = os.WriteFile("substitutedArticle.md", []byte(substitutedArticleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save substituted article to file substitutedArticle.md: %+v", err)
	}

	// answer the questions
	answers := map[string]string{}
	for id, question := range idQuestionMap {
		fmt.Println(id, "-", question)
		system := "Odpowiedz zwięźle na zadane pytanie na podstawie Dokumentu źródłowego dostarczonego w formacie MarkDown"
		user := question + " Dokument źródłowy:\n" + substitutedArticleMD
		answer, err := openai.CompletionCheap(user, system, nil)
		if err != nil {
			log.Printf("OpenAI returned error: %+v", err)
			continue
		}
		fmt.Println(answer)
		fmt.Println()
		answers[id] = answer
	}

	// Verify answer
	result, err := api.VerifyTaskAnswer("arxiv", answers, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
