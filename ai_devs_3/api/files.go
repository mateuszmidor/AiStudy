package api

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func DownloadIfDoesntExistYet(url, destinationFilePath string) error {
	// Check if the destination file already exists
	if _, err := os.Stat(destinationFilePath); err == nil {
		// File already exists, no need to download
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than "file does not exist" occurred
		return errors.Wrap(err, "failed to check if destination file exists")
	}

	// Create the downloads directory if it doesn't exist
	downloadsDir := filepath.Dir(destinationFilePath)
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadsDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "failed to create downloads directory")
		}
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	// Create the destination file to store the downloaded file
	destFile, err := os.Create(destinationFilePath)
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

func UnzipArchive(sourceFilePath, destinationDir string) error {
	// Open the ZIP file
	zipReader, err := zip.OpenReader(sourceFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to open zip file")
	}
	defer zipReader.Close()

	// Extract the files from the ZIP archive
	for _, file := range zipReader.File {
		filePath := filepath.Join(destinationDir, file.Name)

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
