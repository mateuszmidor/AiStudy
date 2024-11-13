package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/joho/godotenv"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

// langfuse-go read following env variables:
// LANGFUSE_HOST
// LANGFUSE_PUBLIC_KEY
// LANGFUSE_SECRET_KEY
func main() {
	maze := `
....B
.###.
...#.
A#...
`
	prompt, err := api.BuildPrompt("prompt.txt", "{{INPUT}}", maze)
	if err != nil {
		log.Fatalf("Error fetching prompt: %+v", err)
	}
	expected := "UP, UP, UP, RIGHT, RIGHT, RIGHT, RIGHT"
	resp, err := tracedCompletionExpert("you are helpful assistant", prompt, "gpt-4o-mini", "json_object", expected)
	if err != nil {
		panic(err)
	}
	if resp.Error == nil {
		fmt.Println(resp.Choices[0].Message.Content)
	} else {
		fmt.Println(resp.Error.Message)
	}
	fmt.Println(maze)
}

func tracedCompletionExpert(system, user, gptModel, format, expectedInResponse string) (*openai.GPTResponse, error) {
	// get completion from OpenAI
	result, err := openai.CompletionExpert(user, system, "", gptModel, format, 1000, 0.0)
	if err != nil {
		panic(err)
	}

	// trace the completion in LangFuse
	l := langfuse.New(context.Background())

	// start new trace that will be visible in LangFuse dashboard
	trace, err := l.Trace(&model.Trace{Name: "maze-solver"})
	if err != nil {
		panic(err)
	}

	// start new span within the trace
	span, err := l.Span(&model.Span{Name: "do-solve", TraceID: trace.ID}, nil)
	if err != nil {
		panic(err)
	}

	// start new generation within the span
	generation, err := l.Generation(
		&model.Generation{
			TraceID: trace.ID,
			Name:    "generation",
			Model:   gptModel,
			ModelParameters: model.M{
				"maxTokens":   "1000",
				"temperature": "0.0",
			},
			Input: []model.M{
				{
					"role":    "system",
					"content": system,
				},
				{
					"role":    "user",
					"content": user,
				},
			},
			Metadata: model.M{
				"difficulty": "medium",
			},
		},
		&span.ID,
	)
	if err != nil {
		panic(err)
	}

	score := 1.0 // score of 0.0 doesnt show up in dashboard, probably SDK problem
	eventMsg := "crashed against the wall!"
	if result.Error != nil {
		generation.Output = model.M{
			"error": result.Error.Message,
		}
	} else {
		if strings.Contains(result.Choices[0].Message.Content, expectedInResponse) {
			score = 6.0
			eventMsg = "escaped from the maze!"
		}
		generation.Output = model.M{
			"completion": result.Choices[0].Message.Content,
		}
		generation.Usage.Input = result.Usage.PromptTokens
		generation.Usage.CompletionTokens = result.Usage.CompletionTokens
		generation.Usage.TotalTokens = result.Usage.TotalTokens
	}

	_, err = l.GenerationEnd(generation)
	if err != nil {
		panic(err)
	}

	// attach event to the generation
	_, err = l.Event(
		&model.Event{
			Name:    eventMsg,
			TraceID: trace.ID,
		},
		&generation.ID,
	)
	if err != nil {
		panic(err)
	}

	// attach score to the trace
	_, err = l.Score(
		&model.Score{
			TraceID: trace.ID,
			Name:    "maze-navigation-score-1..6",
			Value:   score,
		},
	)
	if err != nil {
		panic(err)
	}

	_, err = l.SpanEnd(span)
	if err != nil {
		panic(err)
	}

	l.Flush(context.Background())

	return result, err
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
