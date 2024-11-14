package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	// fetch robot JSON description
	apikey := api.ApiKey()
	taskUrl := "https://centrala.ag3nts.org/data/" + apikey + "/robotid.json"
	robotDescriptionJSON, err := api.FetchData(taskUrl)
	if err != nil {
		log.Fatalf("Error fetching task: %+v", err)
	}
	robotDescription := extractDescription(robotDescriptionJSON)
	fmt.Println("Description")
	fmt.Println(robotDescription)

	// transform description into robotFeatures
	robotFeatures, err := openai.Completion(robotDescription, "From the provided description of a robot, extract all the visual features needed to draw it, respond in english, include just the features", nil, "gpt-4o-mini")
	if err != nil {
		log.Fatalf("Error from openai: %+v", err)
	}
	fmt.Println("Features")
	fmt.Println(robotFeatures)

	// generate image from features
	resp, err := openai.GenerationExpert("Draw a robot:\n"+robotFeatures, "dall-e-3", "1024x1024", "url")
	if err != nil {
		log.Fatalf("Error from openai: %+v", err)
	}
	fmt.Println("URL")
	fmt.Println(resp.Data[0].URL)

	// verify the generated image
	result, err := api.VerifyTaskAnswer("robotid", resp.Data[0].URL, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}

func extractDescription(descriptionJSON string) string {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(descriptionJSON), &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %+v", err)
	}
	description, ok := result["description"].(string)
	if !ok {
		log.Fatalf("Description not found or is not a string")
	}
	return description
}
