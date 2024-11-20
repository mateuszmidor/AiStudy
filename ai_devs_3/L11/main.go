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

	//"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/pkg/errors"
)

const filesZipURL = "https://centrala.ag3nts.org/dane/pliki_z_fabryki.zip"

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

// return filename-keywords map
func getFilenameKeywordsForReports(srcDir string) map[string]string {
	const system = `Zajmujesz się wyodrębnianiem słów kluczowych z tekstu`

	// map of filename -> comma, separated, keywords
	result := map[string]string{}

	// list all files in srcDir
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("failed to read directory: %+v", err)
	}

	// iterate over the files
	for _, file := range files {
		// skip non-txt files
		if strings.ToLower(filepath.Ext(file.Name())) != ".txt" {
			continue
		}

		// if file with extracted keywords already exists - just read it and move on to the next file
		filePath := filepath.Join(srcDir, file.Name())
		keywordsFilePath := filePath + ".keywords"
		if _, err := os.Stat(keywordsFilePath); err == nil {
			content, err := os.ReadFile(keywordsFilePath)
			if err != nil {
				log.Printf("failed to read keywords file %s: %+v\n", keywordsFilePath, err)
			} else {
				result[file.Name()] = string(content)
				continue
			}
		}

		// read the input file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", filePath, err)
			continue
		}

		// extract keywords from content
		log.Print("Extracting keywords for", filePath)
		user, err := api.BuildPrompt("prompt_report.txt", "{{INPUT}}", string(content))
		if err != nil {
			log.Printf("failed to build prompt: %+v\n", err)
			continue
		}
		resultString, err := openai.CompletionCheap(user, system, nil)
		if err != nil {
			log.Printf("openai returend error: %+v", err)
			continue
		}

		// include sector name in keywords, eg "sektor_C2"
		sectorName := strings.Split(file.Name(), "-")[4]
		sectorName = strings.TrimSuffix(sectorName, filepath.Ext(file.Name()))
		onlyName := strings.Split(sectorName, "_")[1] // C2
		keywords := onlyName + ", " + strings.TrimSpace(resultString)
		result[file.Name()] = keywords

		// save keywords to file
		err = os.WriteFile(keywordsFilePath, []byte(keywords), os.ModePerm)
		if err != nil {
			log.Printf("failed to save keywords to file %s: %+v\n", keywordsFilePath, err)
		}
	}
	return result
}

// return person name-keywords map
func getPersonKeywordsForDocuments(srcDir string) map[string]string {
	const system = `Zajmujesz się wyodrębnianiem słów kluczowych z tekstu`

	result := map[string]string{}

	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("failed to read directory: %+v", err)
	}
	for _, file := range files {
		// skip non-txt files
		if strings.ToLower(filepath.Ext(file.Name())) != ".txt" {
			continue
		}

		// if file with extracted keywords already exists - just read it and move on to the next file
		filePath := filepath.Join(srcDir, file.Name())
		keywordsFilePath := filePath + ".keywords"
		if _, err := os.Stat(keywordsFilePath); err == nil {
			content, err := os.ReadFile(keywordsFilePath)
			nameKeywords := strings.Split(string(content), ":")
			if err != nil {
				log.Printf("failed to read keywords file %s: %+v\n", keywordsFilePath, err)
			} else {
				result[nameKeywords[0]] = nameKeywords[1]
				continue
			}
		}

		// read the input file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", filePath, err)
			continue
		}

		// extract the person name and keywords from content
		log.Print("Extracting keywords for", filePath)
		user, err := api.BuildPrompt("prompt_facts.txt", "{{INPUT}}", string(content))
		if err != nil {
			log.Printf("failed to build prompt: %+v\n", err)
			continue
		}
		resultString, err := openai.CompletionCheap(user, system, nil)
		if err != nil {
			log.Fatalf("openai returend error: %+v", err)
		}

		nameKeywords := strings.Split(resultString, ":")
		if len(nameKeywords) != 2 {
			log.Println("invalid name-keywords")
			continue
		}
		result[nameKeywords[0]] = strings.TrimSpace(nameKeywords[1])

		// save person name: keywords to file
		err = os.WriteFile(keywordsFilePath, []byte(resultString), os.ModePerm)
		if err != nil {
			log.Printf("failed to save keywords to file %s: %+v\n", keywordsFilePath, err)
		}
	}
	return result
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

	nameKeywords := getPersonKeywordsForDocuments("downloads/facts")
	// for name, keywords := range nameKeywords {
	// 	fmt.Println(name, ":", keywords)
	// }

	filenameKeywords := getFilenameKeywordsForReports("downloads")
	// for name, keywords := range filenameKeywords {
	// 	fmt.Println(name, ":", keywords)
	// }

	answer := map[string]string{}
	for file, words1 := range filenameKeywords {
		// include keywords from report
		answer[file] = words1

		// add keywords from facts about person if report mentions a person
		for fullName, words2 := range nameKeywords {
			for _, nameSegment := range strings.Split(fullName, " ") { // split eg "Anna Dymna" into [Anna, Dymna]
				if strings.Contains(words1, nameSegment) {
					fmt.Println(file, "mentions", fullName)
					answer[file] = answer[file] + ", " + words2
					break
				}
			}
		}
	}

	// Verify answer
	result, err := api.VerifyTaskAnswer("dokumenty", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
