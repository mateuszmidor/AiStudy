package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func unzipArchive(source, destination string) error {
	// Open the ZIP file
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return errors.Wrap(err, "failed to open zip file")
	}
	defer zipReader.Close()

	// Extract the files from the ZIP archive
	for _, file := range zipReader.File {
		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			// Create directories if necessary
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "failed to create directory")
			}
			continue
		}

		// Create the file
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return errors.Wrap(err, "failed to create file")
		}

		// Open the file inside the ZIP archive
		zipFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return errors.Wrap(err, "failed to open file in zip")
		}

		// Copy the file content
		_, err = io.Copy(destFile, zipFile)
		if err != nil {
			zipFile.Close()
			destFile.Close()
			return errors.Wrap(err, "failed to copy file content")
		}

		// Close the files
		zipFile.Close()
		destFile.Close()
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

// readOrTranscribeImage returns key-value pairs: {filename: transcription} for all .png files found under sourceDir
func readOrTranscribeImage(sourceDir string) (map[string]string, error) {
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

func categorizeFiles(srcDir string) FileCategories {
	// 1. collect all files into one big document
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("failed to read directory: %+v", err)
	}
	var allDocuments string
	for _, file := range files {
		// skip non-txt files
		if strings.ToLower(filepath.Ext(file.Name())) != ".txt" {
			continue
		}
		// transform file.png.txt into file.png
		originalFileName := removeExtraExt(file.Name())

		// read the txt file
		filePath := filepath.Join(srcDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", filePath, err)
			continue
		}

		// compose big document from txt files
		allDocuments += fmt.Sprintf("[%s]\n%s\n\n", originalFileName, string(content))
	}

	// 2. categorize
	system := `
	You will be given a text document composed of a series FileName-FileContent blocks, like this:
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
	user := "The document: \n" + allDocuments
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

func main() {
	// Download the ZIP if doesn't exist yet
	err := downloadZipIfDoesntExistYet(filesZipURL, "downloads/pliki.zip")
	if err != nil {
		log.Fatalf("failed to get the factor files: +%+v", err)
	}

	// Use the existing function to unzip the archive
	err = unzipArchive("downloads/pliki.zip", "downloads/")
	if err != nil {
		log.Fatalf("failed to unzip archive: %+v", err)
	}

	// Transcribe the files
	_, err = readOrTranscribeAudio("downloads/")
	if err != nil {
		log.Fatalf("failed to transcribe audio recordings: %+v", err)
	}

	// Transcribe the files
	_, err = readOrTranscribeImage("downloads/")
	if err != nil {
		log.Fatalf("failed to transcribe audio recordings: %+v", err)
	}

	answer := categorizeFiles("downloads/")
	fmt.Printf("%+v\n", answer)
	result, err := api.VerifyTaskAnswer("kategorie", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
