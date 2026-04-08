package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	// 1. loads "tree.png" image from current directory
	imageData, err := os.ReadFile("tree.png")
	if err != nil {
		log.Fatal(err)
	}

	// 2. initializes langchain for image recognition
	ctx := context.Background()
	llm, err := openai.New(openai.WithModel("gpt-4o-mini"))
	if err != nil {
		log.Fatal(err)
	}

	// 3. uses OpenAI multi-modal model to analyze and recognize the image
	// Encode image to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	imageURL := "data:image/png;base64," + base64Image

	// Create message content with image
	textPart := llms.TextPart("Describe this image in detail")
	imagePart := llms.ImageURLPart(imageURL)

	// Create messages
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

	// 4. prints out description of the image
	fmt.Println("Image Description:")
	fmt.Println(response.Choices[0].Content)
	fmt.Println()

	// 5. prints out token usage
	if len(response.Choices) > 0 && response.Choices[0].GenerationInfo != nil {
		if promptTokens, ok := response.Choices[0].GenerationInfo["PromptTokens"].(int); ok {
			if completionTokens, ok := response.Choices[0].GenerationInfo["CompletionTokens"].(int); ok {
				if totalTokens, ok := response.Choices[0].GenerationInfo["TotalTokens"].(int); ok {
					fmt.Printf("Token Usage: %d input tokens, %d output tokens, %d total tokens\n",
						promptTokens,
						completionTokens,
						totalTokens)
				} else {
					fmt.Println("Token usage information not available")
				}
			} else {
				fmt.Println("Token usage information not available")
			}
		} else {
			fmt.Println("Token usage information not available")
		}
	} else {
		fmt.Println("Token usage information not available")
	}
}
