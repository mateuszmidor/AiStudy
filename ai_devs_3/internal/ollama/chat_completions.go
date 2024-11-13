package ollama

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const chatCompletionsURL = "http://localhost:11434/api/chat"

// https://github.com/ollama/ollama/blob/main/docs/api.md#generate-a-chat-completion
type OllamaRequest struct {
	Model     string    `json:"model"`                // REQUIRED, [llama3:8b, llama3.2-vision:11b, llava:7b]
	Messages  []Message `json:"messages,omitempty"`   // the messages of the chat, this can be used to keep a chat memory
	Tools     any       `json:"tools,omitempty"`      // tools for the model to use if supported. Requires stream to be set to false
	Format    string    `json:"format,omitempty"`     // [json]
	Options   *Options  `json:"options,omitempty"`    // additional model parameters listed in the documentation for the Modelfile such as temperature
	Stream    bool      `json:"stream"`               // default is TRUE if no flag provided
	KeepAlive string    `json:"keep_alive,omitempty"` // controls how long the model will stay loaded into memory following the request (default: 5m)
}

type Message struct {
	Role      string   `json:"role"`                 // the role of the message, either system, user, assistant, or tool
	Content   string   `json:"content"`              // the content of the message
	Images    []string `json:"images,omitempty"`     //  a list of images to include in the message (for multimodal models such as llava)
	ToolCalls []any    `json:"tool_calls,omitempty"` // (optional): a list of tools the model wants to use
}

// https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
type Options struct {
	Seed        int     `json:"seed"`        // e.g. 1.0
	Temperature float64 `json:"temperature"` // e.g. 0
}

type OllamaResponse struct {
	Model              string    `json:"model"` // used model
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"` // reponse message, if success
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"` // prompt tokens
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"` // response tokens
	EvalDuration       int64     `json:"eval_duration"`
	Error              string    `json:"error"` // error message, if error
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: ollama.Completion("Answer in 1 word: what is in the picture", "", ollama.ImageFromFile("avocado.png"), "llama3.2-vision:11b")
func Completion(user, system string, images []string, model string) (string, error) {
	resp, err := CompletionExpert(user, system, images, model, "")
	if err != nil {
		return "", err
	}
	return resp.Message.Content, nil
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
func CompletionExpert(user, system string, images []string, model, responseFormat string) (*OllamaResponse, error) {
	// Initialize the payload
	messages := []Message{}
	if system != "" {
		systemMessage := Message{Role: "system", Content: system}
		messages = append(messages, systemMessage)
	}
	if user != "" || len(images) > 0 {
		userMessage := Message{Role: "user", Content: user, Images: images} // note that llama3.2-vision doesn't support multiple images in single message
		messages = append(messages, userMessage)
	}
	payload := &OllamaRequest{
		Model:    model,
		Stream:   false,
		Messages: messages,
		Format:   responseFormat,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshall ollama request")
	}

	// Create a new request using http.Post
	resp, err := http.Post(chatCompletionsURL, "application/json", bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK || ollamaResponse.Error != "" {
		return nil, errors.Errorf("ollama returned error: %q (http_code %d)", ollamaResponse.Error, resp.StatusCode)
	}

	// SUCCESS
	return &ollamaResponse, nil
}

func ImageFromBytes(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func ImageFromFile(filename string) string {
	imageBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read file %s: %+v\n", filename, err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(imageBytes)
}

func ImageFromURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("failed to fetch image from URL %s: %+v\n", url, err)
		return ""
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read image data from URL %s: %+v\n", url, err)
		return ""
	}

	return ImageFromBytes(imageBytes)
}
