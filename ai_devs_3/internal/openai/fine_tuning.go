package openai

import (
	"encoding/json"
	"log"
	"os"

	"math/rand"
)

// single entry in jsonl:
// {
// 	"messages":[
// 	  {
// 		"role":"system",
// 		"content":"Classify the data into category"
// 	  },
// 	  {
// 		"role":"user",
// 		"content":"tiger"
// 	  },
// 	  {
// 		"role":"assistant",
// 		"content":"CAT"
// 	  }
// 	]
// }

type jsonlMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type jsonlSample struct {
	Messages []jsonlMessage `json:"messages"`
}

type JSONL struct {
	systemMessage   string
	trainingSamples []jsonlSample
}

func NewJSONL(systemMessage string) *JSONL {
	return &JSONL{systemMessage: systemMessage}
}

func (j *JSONL) AddSample(userContent, assistantContent string) {
	sample := jsonlSample{
		Messages: []jsonlMessage{
			{
				Role:    "system",
				Content: j.systemMessage,
			},
			{
				Role:    "user",
				Content: userContent,
			},
			{
				Role:    "assistant",
				Content: assistantContent,
			},
		},
	}
	j.trainingSamples = append(j.trainingSamples, sample)
}

func (j *JSONL) Shuffle() {
	perm := rand.Perm(len(j.trainingSamples))
	shuffledSamples := make([]jsonlSample, len(j.trainingSamples))
	for i, v := range perm {
		shuffledSamples[v] = j.trainingSamples[i]
	}
	j.trainingSamples = shuffledSamples
}

func (j *JSONL) Save(trainingFilename, validationFilename string) error {
	trainingCount := int(float64(len(j.trainingSamples)) * 0.85)
	validationCount := len(j.trainingSamples) - trainingCount
	log.Printf("saving %d samples to %s", trainingCount, trainingFilename)
	log.Printf("saving %d samples to %s", validationCount, validationFilename)

	trainingFile, err := os.Create(trainingFilename)
	if err != nil {
		return err
	}
	defer trainingFile.Close()

	validationFile, err := os.Create(validationFilename)
	if err != nil {
		return err
	}
	defer validationFile.Close()

	trainingEncoder := json.NewEncoder(trainingFile)
	validationEncoder := json.NewEncoder(validationFile)

	for i, sample := range j.trainingSamples {
		if i < trainingCount {
			if err := trainingEncoder.Encode(sample); err != nil {
				return err
			}
		} else {
			if err := validationEncoder.Encode(sample); err != nil {
				return err
			}
		}
	}

	return nil
}
