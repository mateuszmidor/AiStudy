package main

import (
	"fmt"
	"log"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
)

const system = "The polish city which map is attached in fragments has granaries and a fortress"
const user = "Recognize the name of the city in Poland from the provided fragments of the city map. Note: one fragment belongs to a different city that doesn't interes us"

func main() {
	mapFragments := []string{"mapa1.jpg", "mapa2.jpg", "mapa3.jpg", "mapa4.jpg"}
	images := make([]string, len(mapFragments))
	for i, filename := range mapFragments {
		images[i] = ollama.ImageFromFile(filename)
	}
	resp, err := ollama.Completion(user, system, images, "llava:7b")
	if err != nil {
		log.Fatalf("Error from openai: %+v", err)
	}
	fmt.Println(resp)
	fmt.Println()
}
