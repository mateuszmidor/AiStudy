package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
)

const verificationURL = "https://centrala.ag3nts.org/report"

func main() {
	// fetch the task input data
	apikey := os.Getenv("AIDEVS3_API_KEY")
	taskUrl := "https://centrala.ag3nts.org/data/" + apikey + "/cenzura.txt"
	personalData, err := api.FetchData(taskUrl)
	if err != nil {
		log.Fatalf("Error fetching task: %+v", err)
	}

	// fetch the prompt
	prompt, err := api.BuildPrompt("prompt.txt", "{{INPUT}}", personalData)
	if err != nil {
		log.Fatalf("Error fetching prompt: %+v", err)
	}

	// use local LLM to censor personal information
	censoredData, err := ollama.Completion(prompt, "llama3")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Uncensored: %q\n", personalData)
	fmt.Printf("Censored:   %q\n", censoredData)

	// verify the answer
	result, err := api.VerifyTaskAnswer("CENZURA", censoredData, verificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
