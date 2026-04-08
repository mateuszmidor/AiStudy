package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

const imagePrompt = "Describe this image in one sentence"

func main() {
	// 1. loads "tree.png" image from current directory
	imageData, err := os.ReadFile("tree.png")
	if err != nil {
		log.Fatal(err)
	}

	// 2. recognize image using LLM over langchain
	response := recognizeWithOllama(imageData)

	// 3. prints out description of the image
	fmt.Println("Image Description:")
	fmt.Println(response.Choices[0].Content)
	fmt.Println()

	// 4. prints out token usage
	info := response.Choices[0].GenerationInfo
	fmt.Printf("Token Usage: %d input tokens, %d output tokens, %d total tokens\n",
		info["PromptTokens"],
		info["CompletionTokens"],
		info["TotalTokens"])
}

func recognizeWithOllama(imageData []byte) *llms.ContentResponse {
	// prepare the model client
	ctx := context.Background()
	llm, err := ollama.New(ollama.WithModel("qwen3.5:9b"), ollama.WithServerURL("http://localhost:11434")) // vision-capable model
	if err != nil {
		log.Fatal(err)
	}

	// prepare prompt
	textPart := llms.TextPart(imagePrompt)
	imagePart := llms.BinaryPart("image/png", imageData)
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{textPart, imagePart},
		},
	}

	// generate response from the model
	response, err := llm.GenerateContent(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func recognizeWithOpenAI(imageData []byte) *llms.ContentResponse {
	// prepare the model client
	ctx := context.Background()
	llm, err := openai.New(openai.WithModel("gpt-4o-mini")) // vision-capable model
	if err != nil {
		log.Fatal(err)
	}

	// prepare prompt
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	imageURL := "data:image/png;base64," + base64Image
	textPart := llms.TextPart(imagePrompt)
	imagePart := llms.ImageURLPart(imageURL)
	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{textPart, imagePart},
		},
	}

	// Generate response from the model
	response, err := llm.GenerateContent(ctx, messages)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
