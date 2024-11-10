package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/henomis/langfuse-go"
	"github.com/henomis/langfuse-go/model"
	"github.com/joho/godotenv"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = `
<objective>
You are a precise maze solver, who effectively finds a route through given maze from point A to point B, following only allowed fields.
Given a rectangular maze built from ASCII characters, your task is to find a series of allowed moves needed to move from starting point "A" to finish point "B", only stepping on allowed fields marked as ".".
The ASCII characters that build the maze are:
<maze_fields>
. = allowed field, good to step on it
# = obstacle, can not step on it
A = starting point
B = Finish point
</maze_fields>
A move only moves by 1.
The allowed moves for moving around in the maze are:
<allowed_moves>
- UP = move by (X=0, Y=1)
- DOWN = move by (X=0, Y=-1) 
- LEFT = move by (X=-1, Y=0)
- RIGHT = move by (X=1, Y=0)
</allowed_moves>
<example_maze_1>
B.
A.
</example_maze_1>
<example_maze_1_steps>
UP
<example_maze_1_steps>
<example_maze_2>
..
AB
</example_maze_2>
</example_maze_2_steps>
RIGHT
</example_maze_2_steps>
<example_maze_3>
...
...
A#B
</example_maze_3_steps>
UP, RIGHT, RIGHT, DOWN
</example_maze_3_steps>

You should solve the maze by following steps:
1. Translate the maze ASCII representation to coordinates (X,Y) where (0,0) is the top-left corner.
2. Find ALL allowed fields in maze and translate them to coordinates
3. Find positions A and B and translate them to coordinates
4. Find a route through the maze from A to B, it should only go through the allowed fields marked as "."

The expected response is JSON in following format:
{
	"_thoughts": "there is an obstacle on the right so I move up first, ...",
	"steps": "UP, UP, RIGHT, DOWN"
}
`

// langfuse-go read following env variables:
// LANGFUSE_HOST
// LANGFUSE_PUBLIC_KEY
// LANGFUSE_SECRET_KEY
func main() {
	user := `....B
.###.
...#.
A#...
`

	expected := "UP, UP, UP, RIGHT, RIGHT, RIGHT, RIGHT"
	resp, err := tracedCompletionExpert(system, user, "gpt-4o", "json_object", expected)
	if err != nil {
		panic(err)
	}
	if resp.Error == nil {
		fmt.Println(resp.Choices[0].Message.Content)
	} else {
		fmt.Println(resp.Error.Message)
	}
	fmt.Println(user)
}

func tracedCompletionExpert(system, user, gptModel, format, expectedInResponse string) (*openai.GPTResponse, error) {
	// get completion from OpenAI
	result, err := openai.CompletionExpert(system, user, gptModel, format, 1000, 0.0)
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
