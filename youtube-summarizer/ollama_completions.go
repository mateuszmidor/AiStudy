package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

func ollamaGenerateCompletion(prompt string) string {
	fmt.Printf("Prompt:\n%s\n\n", prompt)

	// Initialize the payload
	payload := &OllamaRequest{
		Model:  "llama3",
		Stream: false,
		Prompt: prompt,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	// Specify the URL
	url := "http://localhost:11434/api/generate"

	// Create a new request using http.Post
	fmt.Println("Sending prompt to ollama...")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return ""
	}
	defer resp.Body.Close()
	fmt.Println("Received response from ollama:")

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// Unmarshal the JSON response into an OllamaResponse struct
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	// Assuming you want to return the Response field from the OllamaResponse struct
	return ollamaResponse.Response
}
