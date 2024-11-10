package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type OllamaRequest struct {
	Model  string `json:"model"`            // REQUIRED
	Stream bool   `json:"stream"`           // default is TRUE if no flag provided
	Prompt string `json:"prompt,omitempty"` // user message
	System string `json:"system,omitempty"` // system message
	Format string `json:"format,omitempty"` // [json]
}

type OllamaResponse struct {
	Model              string    `json:"model"` // [llama3:7b]
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

func Completion(prompt, model string) (string, error) {
	resp, err := CompletionExpert("you are a helpful assistant who does as commanded", prompt, model, "")
	if err != nil {
		return "", err
	}
	return resp.Response, nil
}

func CompletionExpert(system, user, model, responseFormat string) (*OllamaResponse, error) {
	// Initialize the payload
	payload := &OllamaRequest{
		Model:  model,
		Stream: false,
		Prompt: user,
		System: system,
		Format: responseFormat,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall ollama request")
	}

	// Specify the URL
	url := "http://localhost:11434/api/generate"

	// Create a new request using http.Post
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to HTTP POST ollama request")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read ollama response")
	}

	// Unmarshal the JSON response into an OllamaResponse struct
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		fmt.Println(string(body))
		return nil, errors.Wrap(err, "failed to unmarshall ollama response")
	}

	// SUCCESS
	return &ollamaResponse, nil
}
