package main

import (
	"fmt"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
)

func main() {
	start := time.Now()
	demo_ollama()
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func demo_ollama() {
	result, err := ollama.Completion("describe what's on the picture", "Be eloquent", ollama.ImageFromFile("../avocado.png"), "llava:7b")
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
