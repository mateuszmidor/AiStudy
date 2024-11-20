package main

import (
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
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/qdrant"

	"github.com/pkg/errors"
	"github.com/yeka/zip"
)

const filesZipURL = "https://centrala.ag3nts.org/dane/pliki_z_fabryki.zip"
const dimensions = 1536
const model = "text-embedding-3-small"
const payloadKeyText = "text" // payload stored and retrieved together with embedding in qdrante

type Embedding = []float64

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

func unzipArchive(source, destination, password string) error {
	// Open the ZIP file
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return errors.Wrap(err, "failed to open zip file")
	}
	defer zipReader.Close()
	// Create the destination directory if it doesn't exist
	downloadsDir := filepath.Dir(destination)
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadsDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed to create destination directory")
		}
	}

	// Extract the files from the ZIP archive
	for _, file := range zipReader.File {
		if file.IsEncrypted() {
			file.SetPassword(password)
		}
		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			// Create directories if necessary
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "failed to create directory")
			}
			continue
		}

		// crate parent directory if doesnt exist
		parentDir := filepath.Dir(filePath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			err = os.MkdirAll(parentDir, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "failed to create parent directory")
			}
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

// return filename-embedding map
func buildEmbeddings(srcDir string) map[string]Embedding {
	result := map[string]Embedding{}

	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("failed to read directory: %+v", err)
	}
	for _, file := range files {
		// skip non-txt files
		if strings.ToLower(filepath.Ext(file.Name())) != ".txt" {
			continue
		}

		// if file with embeddings already exists - just read it and move on to the next file
		filePath := filepath.Join(srcDir, file.Name())
		embeddingsFilePath := filePath + ".embedding"
		if _, err := os.Stat(embeddingsFilePath); err == nil {
			content, err := os.ReadFile(embeddingsFilePath)
			if err != nil {
				log.Printf("failed to read embedding file %s: %+v\n", embeddingsFilePath, err)
			} else {
				embedding := Embedding{}
				err = json.Unmarshal(content, &embedding)
				if err != nil {
					log.Printf("failed to deserialize embedding from file %s: %+v\n", embeddingsFilePath, err)
					continue
				}
				result[file.Name()] = embedding
				continue
			}
		}

		// read the input file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to read file %s: %+v\n", filePath, err)
			continue
		}

		// build embedding
		log.Print("Building embedding for", filePath)
		embedding, err := openai.Embedding(string(content), model, dimensions)
		if err != nil {
			log.Fatalf("openai returend error: %+v", err)
		}

		result[file.Name()] = embedding

		// save embedding to file
		embeddingData, err := json.Marshal(embedding)
		if err != nil {
			log.Printf("failed to serialize embedding for file %s: %+v\n", filePath, err)
			continue
		}

		err = os.WriteFile(embeddingsFilePath, embeddingData, 0644)
		if err != nil {
			log.Printf("failed to write embedding to file %s: %+v\n", embeddingsFilePath, err)
			continue
		}
	}
	return result
}

func filenameToYYYYMMDD(filename string) string {
	dateParts := strings.Split(strings.TrimSuffix(filename, ".txt"), "_")
	var answer string
	if len(dateParts) == 3 {
		year, month, day := dateParts[0], dateParts[1], dateParts[2]
		answer = fmt.Sprintf("%s-%s-%s", year, month, day)
	} else {
		log.Fatal("Filename format is incorrect, expected format: YYYY_MM_DD.txt, got:", filename)
	}
	return answer
}

func main() {
	// Download the ZIP if doesn't exist yet
	err := downloadZipIfDoesntExistYet(filesZipURL, "downloads/pliki.zip")
	if err != nil {
		log.Fatalf("failed to get the factor files: +%+v", err)
	}

	// Use the existing function to unzip the archive
	err = unzipArchive("downloads/pliki.zip", "downloads/", "")
	if err != nil {
		log.Fatalf("failed to unzip archive: %+v", err)
	}

	// Use the existing function to unzip the nested archive
	err = unzipArchive("downloads/weapons_tests.zip", "downloads/weapons_tests", "1670")
	if err != nil {
		log.Fatalf("failed to unzip archive: %+v", err)
	}

	// create embeddings for input txt files
	filenameEmbeddings := buildEmbeddings("downloads/weapons_tests/do-not-share")

	// feed qdrant
	points := []qdrant.Point{}
	for filename, embedding := range filenameEmbeddings {
		payload := map[string]any{payloadKeyText: filename} // keep source filename as "text" entry in payload
		point := qdrant.Point{Vector: embedding, Payload: payload}
		points = append(points, point)
	}
	err = qdrant.FeedDB(points, dimensions)
	if err != nil {
		log.Fatalf("failed to feed qdrant: %+v", err)
	}

	// ask qdrant
	const question = "W raporcie, z którego dnia znajduje się wzmianka o kradzieży prototypu broni?"
	questionEmbedding, err := openai.Embedding(question, model, dimensions)
	if err != nil {
		log.Fatalf("failed to create embedding: %+v", err)
	}
	res, err := qdrant.AskDB(qdrant.Question{Vector: questionEmbedding}, 1)
	if err != nil {
		log.Fatalf("failed to ask qdrant: %+v", err)
	}
	if len(res) == 0 {
		log.Fatal("qdrant returned 0 results")
	}

	// format filename into YYYY-MM-DD date format
	filename := res[0].Payload[payloadKeyText].(string)
	answer := filenameToYYYYMMDD(filename)

	// Verify answer
	result, err := api.VerifyTaskAnswer("wektory", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
