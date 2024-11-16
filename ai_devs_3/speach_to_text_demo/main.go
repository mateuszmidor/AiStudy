package main

import (
	"fmt"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	start := time.Now()
	demo_openai()
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func demo_openai() {
	fmt.Println("demo_openai")
	result, err := openai.SpeachToText("out1.wav")
	if err != nil {
		fmt.Printf("openai error: %+v", err)
		return
	}
	fmt.Println(result)
}
