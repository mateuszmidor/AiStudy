package main

import (
	"fmt"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	start := time.Now()
	// demo_ollama("llava:7b")
	// demo_openai("gpt-4o-mini")
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func demo_ollama() {
	result, err := ollama.Completion("describe what's on the picture", "Be eloquent", ollama.ImageFromFile("../avocado.png"), "llama3.2-vision:11b")
	if err != nil {
		fmt.Printf("ollama error: %+v", err)
		return
	}
	fmt.Println(result)
}

func demo_openai(model string) {
	result, err := openai.Completion("Read ALL text in red found in the picture", "You are OCR expert in reading Polish language", openai.ImageFromURL("https://assets-v2.circle.so/837mal5q2pf3xskhmfuybrh0uwnd"), model)
	if err != nil {
		fmt.Printf("openai error: %+v", err)
		return
	}
	fmt.Println(result)
}
