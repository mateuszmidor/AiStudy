package main

import (
	"context"
	"fmt"
	"log"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/joho/godotenv"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = `
You are an experienced maze solver.
Your task is to find a series of allowed_moves that gets a traveler from position marked as A to position marked as B, avoiding obstacles marked as #.
<allowed_moves>
- UP = move by (0,1), where X axis =0 and Y axis =1
- DOWN = move by (0,-1)
- LEFT = move by (-1,0)
- RIGHT = move by (1,0)
</allowed_moves>
<maze_fields>
. = free field, good to step on it
# = obstacle, can not step on it
A = traveler starting point
B = traveler target point
</maze_fields>
<example_maze_1>
This maze of size 2 x 2 has no obstacles, starting point is at position (0,0), target point is at (0,1):
---
B.
A.
---
The valid series of allowed_moves: UP
</example_maze_1>
<example_maze_2>
This maze of size 2 x 2 has no obstacles, starting point is at position (0,0), target point is at (1,0):
---
..
AB
---
The valid series of allowed_moves: RIGHT
</example_maze_2>
<example_maze_3>
This maze of size 3x3 has one obstacle in the middle, starting point is at position (0,0), target point is at (0,2):
---
...
.#.
A.B
---
The valid series of allowed_moves are: RIGHT, RIGHT
</example_maze_3>
You should solve the maze by following steps:
1. Translate the maze ASCII representation to coordinates (x,y) where (0,0) is the top-left corner
2. Find ALL obstacles in maze and translate them to coordinates
3. Find positions A and B and translate them to coordinates
4. Navigate the maze from A to B avoiding obstacles, step by step, checking for collision with obstacle at every step
`

// langfuse-go read following env variables:
// LANGFUSE_HOST
// LANGFUSE_PUBLIC_KEY
// LANGFUSE_SECRET_KEY
func main() {
	user := `
.#..B	
.###.
...#.
A#...
`
	resp, err := tracedCompletinExpert(system, user, "gpt-4o", "text")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Println(user)
}

func tracedCompletinExpert(system, user, gptModel, format string) (*openai.GPTResponse, error) {
	l := langfuse.New(context.Background())
	result, completionErr := openai.CompletionExpert(system, user, gptModel, "text", 1000, 0.0)

	trace, err := l.Trace(&model.Trace{Name: "maze-solver"})
	if err != nil {
		panic(err)
	}

	span, err := l.Span(&model.Span{Name: "solve", TraceID: trace.ID}, nil)
	if err != nil {
		panic(err)
	}

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
				"key": "value",
			},
		},
		&span.ID,
	)
	if err != nil {
		panic(err)
	}

	if completionErr != nil {
		generation.Output = model.M{
			"error": completionErr,
		}
	} else {
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

	// _, err = l.Score(
	// 	&model.Score{
	// 		TraceID: trace.ID,
	// 		Name:    "test-score",
	// 		Value:   0.9,
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = l.Event(
	// 	&model.Event{
	// 		Name:    "test-event",
	// 		TraceID: trace.ID,
	// 		Metadata: model.M{
	// 			"key": "value",
	// 		},
	// 		Input: model.M{
	// 			"key": "value",
	// 		},
	// 		Output: model.M{
	// 			"key": "value",
	// 		},
	// 	},
	// 	&generation.ID,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	span.Output = "Finished span"
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
