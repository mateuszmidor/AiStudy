package main

import (
	"fmt"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	start := time.Now()
	demo_ollama("llava:7b")
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)

	fmt.Println()

	start = time.Now()
	demo_openai("gpt-4o-mini")
	elapsed = time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

// models tested on Asus TUF: Ryzen4800H 16 cores + GTX1660Ti 6GB VRAM + 16GB RAM
// - llava:7b - 5sec (GPU)
// - llama3.2-vision:11b - 18s (CPU, needs 8GB VRAM)
// models tested on asus MacBook PRO 16-inch 2019: Intel i7 2.6GHz 6 cores + integrated graphics + 16GB RAM. All run on CPU
// - llava:7b - 12sec
// - llama3.2-vision:11b - 4m13s
func demo_ollama(model string) {
	fmt.Println("demo_ollama", model)
	result, err := ollama.Completion("describe what's on the picture", "Be eloquent", ollama.ImageFromFile("../avocado.png"), model)
	if err != nil {
		fmt.Printf("ollama error: %+v", err)
		return
	}
	fmt.Println(result)
}

func demo_openai(model string) {
	fmt.Println("demo_openai", model)
	result, err := openai.Completion("describe what's on the picture", "Be eloquent", openai.ImageFromFile("../avocado.png"), model)
	if err != nil {
		fmt.Printf("openai error: %+v", err)
		return
	}
	fmt.Println(result)
}
