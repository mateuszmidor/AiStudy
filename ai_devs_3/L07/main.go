package main

import (
	"fmt"
	"log"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = "You are an expert in maps and you know cities in Poland very well. Your task is to find out polish city from street names."
const user = `
1. Recognize the street names from attached map fragments. 
2. Find out the name of the city in Poland. You are looking for polish city that has granaries and a fortress. 
Note: one map fragment belongs to a different city that doesn't interests us.
`

func main() {
	filenames := []string{"mapa1.jpg", "mapa2.jpg", "mapa3.jpg", "mapa4.jpg"}
	images := make([]string, len(filenames))
	for i, filename := range filenames {
		images[i] = openai.ImageFromFile(filename)
	}
	resp, err := openai.Completion(user, system, images, "gpt-4o")
	if err != nil {
		log.Fatalf("Error from openai: %+v", err)
	}
	fmt.Println(resp)
}
