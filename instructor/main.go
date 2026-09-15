package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jxnl/instructor-go/pkg/instructor"
	openai "github.com/sashabaranov/go-openai"
)

// json + jsonschema tags are enough for Instructor to build the JSON schema:
// it injects the schema into the prompt, parses the answer and validates it.
type City struct {
	City    string `json:"city" jsonschema:"title=city,description=Name of a city"`
	Country string `json:"country" jsonschema:"title=country,description=Country the city is located in"`
}

func main() {
	ctx := context.Background()

	client := instructor.FromOpenAI(
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		instructor.WithMode(instructor.ModeJSON),
		instructor.WithMaxRetries(3),
	)

	// parity with the baml demo: hardcoded country
	const country = "China"

	var city City
	_, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Pick one random city in %s.", country),
			},
		},
	}, &city)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Random city: %s, %s", city.City, city.Country)
}
