package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	openAIURL = "https://api.openai.com/v1/chat/completions"
)

type GPTRequest struct {
	Model          string    `json:"model"`                     // REQUIRED, [gpt-3.5-turbo, gpt-4-turbo-preview, gpt-4o-mini, gpt-4o]
	Messages       []Message `json:"messages"`                  // REQUIRED, at least 1 "user" message
	ResponseFormat *Format   `json:"response_format,omitempty"` // [text, json_object]
	NumAnswers     uint      `json:"n,omitempty"`               // [1..+oo], default: 1; cheapest option
	MaxTokens      int       `json:"max_tokens,omitempty"`      // [1..+oo], default: ?; max tokens generated for answer before the generation is hard-cut
	Temperature    float32   `json:"temperature,omitempty"`     // [0.0..2.0], default: 0 (auto-select); use high for creativity and randomness
}

type Message struct {
	Role    string `json:"role"` // [user, system, assistant]; assistant means a previous GPT response; include it for interaction continuity
	Content string `json:"content"`
}

type Format struct {
	Type string `json:"type"` // [text, json_object]; if json_object -> MUST ask gpt directly to respond in JSON format
}

type GPTResponse struct {
	Choices []Choice  `json:"choices"` // number of returned choices is directly related to GPTRequest.NumAnsers
	Created int64     `json:"created"`
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Object  string    `json:"object"`
	Usage   Usage     `json:"usage"`
	Error   *GPTError `json:"error"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"` // [stop, length, content_filter]; stop means natural stop while length means MaxTokens hit
	Index        int     `json:"index"`         // index in the list of choices
	Message      Message `json:"message"`
	Logprobs     *string `json:"logprobs"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GPTError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

func (e *GPTError) Error() string {
	return e.Message
}

func completion(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY") // Get the API key from the environment variable
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set")
	}

	reqBody := GPTRequest{
		Model:          "gpt-4o-mini",
		ResponseFormat: &Format{Type: "text"},
		NumAnswers:     1,
		MaxTokens:      256,
		Temperature:    0.0,
		Messages: []Message{
			{
				Role: "system", Content: "do your best to help the user by answering concisely and precisely to user's question",
			},
			{
				Role: "user", Content: prompt,
			},
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openAIURL, bytes.NewReader(reqBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var gptResp GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return "", err
	}

	if len(gptResp.Choices) > 0 {
		if len(gptResp.Choices) > 1 {
			fmt.Println(gptResp.Choices[1:])
		}
		return gptResp.Choices[0].Message.Content, nil
	}

	return "", gptResp.Error
}
