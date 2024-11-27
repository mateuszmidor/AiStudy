package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = "classify the numbers"
const correct = "correct"
const incorrect = "wrong"

func main() {
	// download task data
	err := api.DownloadIfDoesntExistYet("https://centrala.ag3nts.org/dane/lab_data.zip", "downloads/lab_data.zip")
	if err != nil {
		log.Fatal(err)
	}

	// extract task data
	err = api.UnzipArchive("downloads/lab_data.zip", "downloads/")
	if err != nil {
		log.Fatal(err)
	}

	// create fine tuning file for openai
	prepareTrainingJSONL("downloads/correct.txt", "downloads/incorrect.txt")

	// classify input
	answer := classify("downloads/verify.txt", "ft:gpt-4o-mini-2024-07-18:personal:aidev3-s04-e02-v5:AYJsQoYD")

	// verify the generated image
	fmt.Println("answer:", answer)
	result, err := api.VerifyTaskAnswer("research", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}

func prepareTrainingJSONL(correctSamplesFilename, incorrectSamplesFilename string) {
	jsonl := openai.NewJSONL(system)
	addLinesAsSamples(jsonl, correctSamplesFilename, correct)
	addLinesAsSamples(jsonl, incorrectSamplesFilename, incorrect)
	jsonl.Shuffle()
	jsonl.Save("trainingData.jsonl", "validationData.jsonl")
}

func addLinesAsSamples(jsonl *openai.JSONL, filename, category string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		jsonl.AddSample("the numbers: "+line, category)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func classify(filename, finetunedModel string) []string {
	correctIDs := []string{}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			// good
		} else {
			log.Printf("Unexpected format: %s", line)
			continue
		}
		id, data := parts[0], parts[1]
		rsp, err := openai.Completion("the numbers: "+data, system, nil, finetunedModel)
		if err != nil {
			log.Println("openai failed:", err.Error())
		}
		fmt.Printf("%s [%s] is %s\n", id, data, rsp)
		if rsp == correct {
			correctIDs = append(correctIDs, id)
		}
	}
	return correctIDs
}
