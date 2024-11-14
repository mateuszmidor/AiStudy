package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const openAIURL = "https://api.openai.com/v1/images/generations"

type DalleGenerateRequest struct {
	Prompt string `json:"prompt"`                    // REQUIRED
	Model  string `json:"model,omitempty"`           // [dall-e-2,dall-e-3], default: dall-e-2
	N      int    `json:"n,omitempty"`               // num images to generate, default: 1
	Size   string `json:"size,omitempty"`            // [256x256, 512x512, or 1024x1024] for dall-e-2, [1024x1024, 1792x1024, or 1024x1792] for dall-e-3, default: 1024x1024
	Format string `json:"response_format,omitempty"` // [url, b64_json], default: url
}

type DalleGenerateResponse struct {
	Created int64 `json:"created"` // unix timestamp
	Data    []struct {
		URL         string `json:"url"`      // depending on DalleGenerateRequest.Format: either this
		Base64Image string `json:"b64_json"` // or this
	} `json:"data"`
}

type DalleErrorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// max 1000 characters for dall-e-2 and 4000 characters for dall-e-3
func GenerationExpert(prompt, model, size, format string) (*DalleGenerateResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY") // Get the API key from the environment variable
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is not set")
	}

	// send request
	req, err := preparePromptToImageRequest(prompt, model, size, format, apiKey)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var errorResponse DalleErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, err
		}
		errorMsg := fmt.Sprintf("request failed [code:%d]: %s [type:%s]", resp.StatusCode, errorResponse.Error.Message, errorResponse.Error.Type)
		return nil, errors.New(errorMsg)
	}

	// deserialize response
	var generateResponse DalleGenerateResponse
	err = json.Unmarshal(body, &generateResponse)
	if err != nil {
		return nil, err
	}

	// assuming you want to return the URL of the first generated image
	return &generateResponse, nil
}

// https://platform.openai.com/docs/api-reference/images/create
func preparePromptToImageRequest(prompt, model, size, format string, apiKey string) (*http.Request, error) {
	// prepare request body
	reqBody := DalleGenerateRequest{
		Prompt: prompt,
		Model:  model,
		Size:   size,
		Format: format,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	req, err := http.NewRequest("POST", openAIURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return req, nil
}
