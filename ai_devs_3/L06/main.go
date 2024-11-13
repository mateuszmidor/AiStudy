package main

import (
	"archive/zip"
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

const verificationURL = "https://centrala.ag3nts.org/report"
const recordingsZipURL = "https://centrala.ag3nts.org/dane/przesluchania.zip"

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

// readOrTranscribe returns key-value pairs: {filename: transcription} for all .m4a files found under sourceDir
func readOrTranscribe(sourceDir string) (map[string]string, error) {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	transcriptions := make(map[string]string)

	for _, file := range files {
		// skip non-audio files
		if strings.ToLower(filepath.Ext(file.Name())) != ".m4a" {
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

func main() {
	// Download the ZIP if doesn't exist yet
	err := downloadZipIfDoesntExistYet(recordingsZipURL, "downloads/przesluchania.zip")
	if err != nil {
		log.Fatalf("failed to get the audio recordings: +%+v", err)
	}

	// Use the existing function to unzip the archive
	err = unzipArchive("downloads/przesluchania.zip", "downloads/")
	if err != nil {
		log.Fatalf("failed to unzip archive: %+v", err)
	}

	// Transcribe the files
	transcriptions, err := readOrTranscribe("downloads/")
	if err != nil {
		log.Fatalf("failed to transcribe audio recordings: %+v", err)
	}

	// collect transcriptions into single block of text
	allTranscriptions := ""
	for sourceFile, transcription := range transcriptions {
		allTranscriptions = allTranscriptions + sourceFile + "\n" + transcription + "\n"
	}
	fmt.Println(allTranscriptions)

	// build prompt
	prompt, err := api.BuildPrompt("prompt.txt", "{{INPUT}}", allTranscriptions)
	if err != nil {
		log.Fatalf("failed to build prompt: %+v", err)
	}

	fmt.Println(prompt)
	// ask LLM
	answer, err := openai.CompletionCheap(prompt, "", "")
	if err != nil {
		log.Fatalf("Error getting answer for question '%s': %v\n", prompt, err)
	}
	fmt.Println("answer:\n" + answer)

	// verify the answer
	result, err := api.VerifyTaskAnswer("mp3", answer, verificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
