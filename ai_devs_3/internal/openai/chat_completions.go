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
	ResponseFormat *Format          `json:"response_format,omitempty"` // [text, json_object, json_schema]
	NumAnswers     uint             `json:"n,omitempty"`               // [1..+oo], default: 1; cheapest option
	MaxTokens      int              `json:"max_tokens,omitempty"`      // [1..+oo], default: ?; max tokens generated for answer before the generation is hard-cut
	Temperature    float32          `json:"temperature,omitempty"`     // [0.0..2.0], default: 0 (auto-select); use high for creativity and randomness
	Tools          []Tool           `json:"tools,omitempty"`           // list of functions available for the model to call
	ToolChoice     string           `json:"tool_choice,omitempty"`     // [none,auto], default: auto, none forces GPT to use no tools (functions)
}

type RequestMessage struct {
	Role       string        `json:"role"`                 // [user, system, assistant, tool]; assistant means a previous GPT response; include it for interaction continuity
	Content    []ContentItem `json:"content"`              // list of content items; OpenAI supports sending multiple images in single message
	ToolCalls  []ToolCall    `json:"tool_calls,omitempty"` // or GPT function call (response)
	ToolCallID string        `json:"tool_call_id"`         // correlation ID for relevant ToolCall
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

type Tool struct {
	Type     string   `json:"type"`     // REQUIRED, [function]
	Function Function `json:"function"` // REQUIRED
}

type Function struct {
	Name        string      `json:"name"` // REQUIRED
	Description string      `json:"description,omitempty"`
	Parameters  *Parameters `json:"parameters,omitempty"`
}

// JSON Schema style
type Parameters struct {
	Type       string              `json:"type"`       // REQUIRED, [object]
	Properties map[string]Property `json:"properties"` // map of property-name:property-details
}

type Property struct {
	Type        string `json:"type"`                  // REQUIRED, [string, integer]
	Description string `json:"description,omitempty"` // what is the purpose of this property? GPT likes details like this
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
	FinishReason string          `json:"finish_reason"` // [stop, length, content_filter, tool_calls]; stop means natural stop while length means MaxTokens hit
	Index        int             `json:"index"`         // index in the list of choices
	Message      ResponseMessage `json:"message"`
	Logprobs     *string         `json:"logprobs"`
}

type ResponseMessage struct {
	Role      string     `json:"role"` // [assistant]; assistant means a previous GPT response; include it for interaction continuity
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"` // or GPT function call (request)
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
	Type     string           `json:"type"`
	ID       string           `json:"id"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`      // function to call, e.g. "get_temperature_at_location_in_celsius"
	Arguments string `json:"arguments"` // arguments for the function as JSON, e.g. {"location": "Gdańsk"}
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

// ChatWithMemory keeps discussion history
type ChatWithMemory struct {
	apiKey       string
	systemPrompt string
	model        string
	maxTokens    int
	messages     []RequestMessage
	debug        bool
}

func NewChatWithMemory(system, model string, maxTokens int, debug bool) (*ChatWithMemory, error) {

	// Get the API key from the environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is not set")
	}

	chat := &ChatWithMemory{
		apiKey:       apiKey,
		systemPrompt: system,
		model:        model,
		maxTokens:    maxTokens,
		messages:     []RequestMessage{},
		debug:        debug,
	}

	// first, attach system message if provided
	if system != "" {
		systemTextContent := ContentItem{Type: "text", Text: system}
		systemMessage := RequestMessage{Role: "system", Content: []ContentItem{systemTextContent}}
		chat.pushMessage(systemMessage)
	}

	return chat, nil
}
func (c *ChatWithMemory) pushMessage(msg RequestMessage) {
	if c.debug {
		fmt.Println(FormatMessage(msg))
	}

	// "tool" request
	if msg.Role == "tool" {
		// "tool" response must immediately follow "tool" request
		for i, message := range c.messages {
			if message.ToolCallID == msg.ToolCallID {
				log.Println("znaleziono tool request na pozycji:", i)
				c.messages = append(c.messages[:i+1], append([]RequestMessage{msg}, c.messages[i+1:]...)...)
				break
			}
		}
	}

	// regular message
	c.messages = append(c.messages, msg)
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: User("describe what's on the picture", openai.ImageFromFile("./avocado.png"), nil, "text", 0.0)
func (c *ChatWithMemory) User(userPrompt string, images []string, tools []Tool, responseFormat string, temperature float32) (*GPTResponse, error) {
	// attach user messages, if provided
	content := []ContentItem{}
	if userPrompt != "" {
		userTextContent := ContentItem{Type: "text", Text: userPrompt}
		content = append(content, userTextContent)
	}
	for _, image := range images {
		userImageContent := ContentItem{Type: "image_url", ImageURL: &ImageURL{URL: image}}
		content = append(content, userImageContent)
	}

	// newMessages := make([]RequestMessage, len(c.messages))
	// copy(newMessages, c.messages)
	if len(content) > 0 {
		userMessage := RequestMessage{Role: "user", Content: content}
		// newMessages = append(newMessages, userMessage)
		c.pushMessage(userMessage)
	}

	toolChoice := ""
	if len(tools) > 0 {
		toolChoice = "auto"
	}
	reqBody := GPTRequest{
		Model:          c.model,
		ResponseFormat: &Format{Type: responseFormat},
		NumAnswers:     1,
		MaxTokens:      c.maxTokens,
		Temperature:    temperature,
		Messages:       c.messages,
		Tools:          tools,
		ToolChoice:     toolChoice,
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
	if gptResp.Error != nil {
		return nil, errors.New(gptResp.Error.Error())
	}
	if len(gptResp.Choices) == 0 {
		return nil, errors.New("llm returned 0 responses")
	}

	// include response in conversation history
	for _, choice := range gptResp.Choices {
		assistantMsg := RequestMessage{
			Role:      choice.Message.Role,
			Content:   []ContentItem{{Type: "text", Text: choice.Message.Content}},
			ToolCalls: choice.Message.ToolCalls,
		}
		// newMessages = append(newMessages, msg)
		c.pushMessage(assistantMsg)
	}
	// c.messages = newMessages
	if Debug {
		log.Printf("Usage: %+v", gptResp.Usage)
	}
	return &gptResp, nil
}

func (c *ChatWithMemory) ToolResponse(response string, toolCallID string) (*GPTResponse, error) {
	// newMessages := make([]RequestMessage, len(c.messages))
	// copy(newMessages, c.messages)
	toolMessage := RequestMessage{Role: "tool", Content: []ContentItem{{Type: "text", Text: response}}, ToolCallID: toolCallID}
	// newMessages = append(newMessages, toolMessage)
	c.pushMessage(toolMessage)

	reqBody := GPTRequest{
		Model:      c.model,
		NumAnswers: 1,
		MaxTokens:  c.maxTokens,
		Messages:   c.messages,
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
	if gptResp.Error != nil {
		return nil, errors.New(gptResp.Error.Error())
	}
	if len(gptResp.Choices) == 0 {
		return nil, errors.New("llm returned 0 responses")
	}

	// include response in conversation history
	for _, choice := range gptResp.Choices {
		msg := RequestMessage{
			Role:      choice.Message.Role,
			Content:   []ContentItem{{Type: "text", Text: choice.Message.Content}},
			ToolCalls: choice.Message.ToolCalls,
		}
		// newMessages = append(newMessages, msg)
		c.pushMessage(msg)
	}
	// c.messages = newMessages
	if Debug {
		log.Printf("Usage: %+v", gptResp.Usage)
	}
	return &gptResp, nil
}

func (c *ChatWithMemory) DumpConversation() string {
	var conversationHistory string
	for i, message := range c.messages {
		conversationHistory += fmt.Sprintf("%d.\n", i)
		conversationHistory += FormatMessage(message)
	}
	return conversationHistory
}

func FormatMessage(message RequestMessage) string {
	messageStr := ""

	for _, toolCall := range message.ToolCalls {
		messageStr += fmt.Sprintf("%s: [Tool Call %s] %s with arguments %s\n", message.Role, toolCall.ID, toolCall.Function.Name, toolCall.Function.Arguments)
	}

	for _, contentItem := range message.Content {
		if contentItem.Text == "" {
			continue
		}
		if contentItem.Type == "text" {
			messageStr += fmt.Sprintf("%s: %s\n", message.Role, contentItem.Text)
		} else if contentItem.Type == "image_url" && contentItem.ImageURL != nil {
			messageStr += fmt.Sprintf("%s: [Image] %s\n", message.Role, contentItem.ImageURL.URL)
		}
	}
	return messageStr
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: openai.CompletionCheap("describe what's on the picture", "Use max 3 words", openai.ImageFromFile("./avocado.png"))
func CompletionCheap(user, system string, images []string) (string, error) {
	return Completion(user, system, images, "gpt-4o-mini")
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: openai.CompletionStrong("describe what's on the picture", "Use max 3 words", openai.ImageFromFile("./avocado.png"))
func CompletionStrong(user, system string, images []string) (string, error) {
	return Completion(user, system, images, "gpt-4o")
}

// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: openai.Completion("describe what's on the picture", "Use max 3 words", openai.ImageFromFile("./avocado.png"), "gpt-4o-mini")
func Completion(user, system string, images []string, model string) (string, error) {
	gptResp, err := CompletionExpert(user, system, images, nil, model, "text", 1000, 0.0)
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
// For image param, use function: ImageFromBytes, ImageFromFile, ImageFromURL.
// Example: openai.Completion("describe what's on the picture", "Use max 3 words", openai.ImageFromFile("./avocado.png"), nil, "gpt-4o-mini", "text", 256, 0.0)
func CompletionExpert(user, system string, images []string, tools []Tool, model, responseFormat string, maxTokens int, temperature float32) (*GPTResponse, error) {
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
	for _, image := range images {
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
		Tools:          tools,
		ToolChoice:     "auto",
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

	if Debug {
		log.Printf("Usage: %+v", gptResp.Usage)
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
