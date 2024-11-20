package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const embeddingsURL = "https://api.openai.com/v1/embeddings"

type EmbeddingRequest struct {
	Model      string `json:"model"`                // [text-embedding-3-small, text-embedding-3-large, text-embedding-ada-002]
	Input      string `json:"input"`                // text to create embedding for
	Dimensions int    `json:"dimensions,omitempty"` // e.g. 1536, 3072
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func Embedding(input, model string, dimensions int) ([]float64, error) {
	resp, err := EmbeddingExpert(input, model, dimensions)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("embedding returned 0 results")
	}
	return resp.Data[0].Embedding, nil
}

func EmbeddingExpert(input, model string, dimensions int) (*EmbeddingResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY") // Get the API key from the environment variable
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is not set")
	}

	reqBody := EmbeddingRequest{
		Model:      model,
		Input:      input,
		Dimensions: dimensions,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	req, err := http.NewRequest("POST", embeddingsURL, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var gptResp EmbeddingResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return nil, errors.New(err.Error())
	}

	if Debug {
		log.Printf("Usage: %+v", gptResp.Usage)
	}
	return &gptResp, nil
}
