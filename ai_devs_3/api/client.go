package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const VerificationURL = "https://centrala.ag3nts.org/report"

// apiRequest represents the structure for API requests
type apiRequest struct {
	Task   string      `json:"task"`
	APIKey string      `json:"apikey"`
	Answer interface{} `json:"answer"`
}

// apiResponse represents the structure for API responses
type apiResponse struct {
	Code    int    `json:"code"`    // negative code means error
	Message string `json:"message"` // response message
}

func ApiKey() string {
	key := os.Getenv("AIDEVS3_API_KEY")
	if key == "" {
		log.Fatal("AIDEVS3_API_KEY is not set")
	}
	return key
}

// VerifyTaskAnswer sends task answer to the API and returns error and response message
func VerifyTaskAnswer(task string, answer interface{}, verificationURL string) (string, error) {
	client := &http.Client{}

	// Create the request body
	reqBody := apiRequest{
		Task:   task,
		APIKey: ApiKey(),
		Answer: answer,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", errors.Wrap(err, "error marshaling request body")
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", verificationURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", errors.Wrap(err, "error creating request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error sending request")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	// Unmarshal the response
	var apiResp apiResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return "", errors.Wrap(err, "error unmarshaling response")
	}

	// Check for API error
	if apiResp.Code < 0 {
		return "", errors.Errorf("API error: %s (code: %d)", apiResp.Message, apiResp.Code)
	}

	return apiResp.Message, nil
}

func FetchData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching resource from %q", url)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "error reading http response body fetched from %q", url)
	}
	return string(data), nil
}

func BuildPrompt(filename, placeholder, data string) (string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "error reading file %q", filename)
	}

	// Replace the placeholder with the provided data
	updatedContent := strings.ReplaceAll(string(fileContent), placeholder, data)
	return updatedContent, nil
}
