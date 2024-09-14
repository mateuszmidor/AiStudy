package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
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
	PromptEvalCount    int       `json:"prompt_eval_count"` // prompt tokens
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"` // response tokens
	EvalDuration       int64     `json:"eval_duration"`
}

func OllamaGenerateCompletion(prompt string) string {
	slog.Info("Prompt: \n" + prompt)

	// Initialize the payload
	payload := &OllamaRequest{
		Model:  "llama3",
		Stream: false,
		Prompt: prompt,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Error marshaling JSON", "error", err)
		return ""
	}

	// Specify the URL
	url := "http://localhost:11434/api/generate"

	// Create a new request using http.Post
	slog.Info("Sending prompt to ollama...")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error sending POST request", "error", err)
		return ""
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
		return ""
	}

	// Unmarshal the JSON response into an OllamaResponse struct
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		slog.Error("Error unmarshaling JSON:", "error", err)
		return ""
	}

	slog.Debug("Received response from ollama:")
	slog.Debug("- input tokens:", "count", ollamaResponse.PromptEvalCount)
	slog.Debug("- output tokens:", "count", ollamaResponse.EvalCount)
	return ollamaResponse.Response
}
