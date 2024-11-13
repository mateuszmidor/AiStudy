package openai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const (
	chatCompletionsURL = "https://api.openai.com/v1/chat/completions"
	base64ImagePrefix  = "data:image/jpeg;base64," // used to send image data in request
)

type GPTRequest struct {
	Model          string           `json:"model"`                     // REQUIRED, [gpt-3.5-turbo, gpt-4-turbo-preview, gpt-4o-mini, gpt-4o]
	Messages       []RequestMessage `json:"messages"`                  // REQUIRED, at least 1 "user" message
	ResponseFormat *Format          `json:"response_format,omitempty"` // [text, json_object]
	NumAnswers     uint             `json:"n,omitempty"`               // [1..+oo], default: 1; cheapest option
	MaxTokens      int              `json:"max_tokens,omitempty"`      // [1..+oo], default: ?; max tokens generated for answer before the generation is hard-cut
	Temperature    float32          `json:"temperature,omitempty"`     // [0.0..2.0], default: 0 (auto-select); use high for creativity and randomness
}

type RequestMessage struct {
	Role    string        `json:"role"`    // [user, system, assistant]; assistant means a previous GPT response; include it for interaction continuity
	Content []ContentItem `json:"content"` // list of content items
}

type ContentItem struct {
	Type     string    `json:"type"`                // [text, image_url]
	Text     string    `json:"text,omitempty"`      // text content, only present if Type is "text"
	ImageURL *ImageURL `json:"image_url,omitempty"` // image URL content, only present if Type is "image_url"
}

type ImageURL struct {
	URL    string `json:"url"`              // data or URL of the image, formatted as "data:image/jpeg;base64,{base64_image_data}" or "https://images.com/avocado.png"
	Detail string `json:"detail,omitempty"` // should the model interpret low or high resolution image [low, high, auto], default: auto. Low costs always 85 tokens. High costs 170 x num of 512x512px squares that build the image +85
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
	FinishReason string          `json:"finish_reason"` // [stop, length, content_filter]; stop means natural stop while length means MaxTokens hit
	Index        int             `json:"index"`         // index in the list of choices
	Message      ResponseMessage `json:"message"`
	Logprobs     *string         `json:"logprobs"`
}

type ResponseMessage struct {
	Role    string `json:"role"` // [assistant]; assistant means a previous GPT response; include it for interaction continuity
	Content string `json:"content"`
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

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL
func CompletionCheap(user, system, image string) (string, error) {
	return Completion(user, system, image, "gpt-4o-mini")
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL
func CompletionStrong(user, system, image string) (string, error) {
	return Completion(user, system, image, "gpt-4o")
}

func Completion(user, system, image, model string) (string, error) {
	gptResp, err := CompletionExpert(system, user, image, model, "text", 1000, 0.0)
	if err != nil {
		return "", err
	}
	if len(gptResp.Choices) > 0 {
		if len(gptResp.Choices) > 1 {
			fmt.Println(gptResp.Choices[1:])
		}
		return gptResp.Choices[0].Message.Content, nil
	}

	return "", errors.New(gptResp.Error.Error())
}

// CompletionExpert generates chat completion, with image support.
// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL
func CompletionExpert(system, user, image, model, responseFormat string, maxTokens int, temperature float32) (*GPTResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY") // Get the API key from the environment variable
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is not set")
	}

	// collect all messages that we want to send to OpenAI
	messages := []RequestMessage{}

	// first, attach system message if provided
	if system != "" {
		systemTextContent := ContentItem{Type: "text", Text: system}
		systemMessage := RequestMessage{Role: "system", Content: []ContentItem{systemTextContent}}
		messages = append(messages, systemMessage)
	}

	// then, attach user messages, if provided
	content := []ContentItem{}
	if user != "" {
		userTextContent := ContentItem{Type: "text", Text: user}
		content = append(content, userTextContent)
	}
	if image != "" {
		userImageContent := ContentItem{Type: "image_url", ImageURL: &ImageURL{URL: image}}
		content = append(content, userImageContent)
	}
	if len(content) > 0 {
		userMessage := RequestMessage{Role: "user", Content: content}
		messages = append(messages, userMessage)
	}

	reqBody := GPTRequest{
		Model:          model,
		ResponseFormat: &Format{Type: responseFormat},
		NumAnswers:     1,
		MaxTokens:      maxTokens,
		Temperature:    temperature,
		Messages:       messages,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	req, err := http.NewRequest("POST", chatCompletionsURL, bytes.NewReader(reqBytes))
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

	var gptResp GPTResponse
	if err := json.Unmarshal(body, &gptResp); err != nil {
		return nil, errors.New(err.Error())
	}

	return &gptResp, nil
}

func ImageFromBytes(bytes []byte) string {
	return base64ImagePrefix + base64.StdEncoding.EncodeToString(bytes)
}

func ImageFromFile(filename string) string {
	imageBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read file %s: %+v\n", filename, err)
		return ""
	}

	return base64ImagePrefix + base64.StdEncoding.EncodeToString(imageBytes)
}

func ImageFromURL(url string) string {
	return url // yes, that simple
}
