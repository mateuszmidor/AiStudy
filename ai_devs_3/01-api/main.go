package main

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
)

const dataURL = "https://poligon.aidevs.pl/dane.txt"
const verifyURL = "https://poligon.aidevs.pl/verify"
const taskName = "POLIGON"

func main() {
	logger := slog.Default()

	// Fetch data from dataURL
	resp, err := http.Get(dataURL)
	if err != nil {
		logger.Error("Error fetching data", "error", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response", "error", err)
		return
	}

	// Convert the data to a string and remove any leading or trailing whitespace
	dataString := strings.TrimSpace(string(data))

	// Split dataString by newline into a slice of strings
	lines := strings.Split(dataString, "\n")

	// Post the answer for verification
	err, msg := api.PostAnswer(taskName, lines, verifyURL)
	if err != nil {
		logger.Error("Error posting answer", "error", err)
		return
	}
	logger.Info("Verification response", "message", msg)
}
