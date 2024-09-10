package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// APIRequest represents the structure for API requests
type APIRequest struct {
	Task   string      `json:"task"`
	APIKey string      `json:"apikey"`
	Answer interface{} `json:"answer"`
}

// APIResponse represents the structure for API responses
type APIResponse struct {
	Code    int    `json:"code"`    // negative code means error
	Message string `json:"message"` // response message
}

// PostAnswer sends task answer to the API and returns error and response message
func PostAnswer(task string, answer interface{}, address string) (error, string) {
	client := &http.Client{}

	// Create the request body
	reqBody := APIRequest{
		Task:   task,
		APIKey: os.Getenv("AIDEVS3_API_KEY"),
		Answer: answer,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %w", err), ""
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err), ""
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err), ""
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err), ""
	}

	// Unmarshal the response
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err), ""
	}

	// Check for API error
	if apiResp.Code < 0 {
		return fmt.Errorf("API error: %s (code: %d)", apiResp.Message, apiResp.Code), ""
	}

	return nil, apiResp.Message
}
