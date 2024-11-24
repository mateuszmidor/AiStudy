package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const barbaraURL = "https://centrala.ag3nts.org/dane/barbara.txt"
const peopleURL = "https://centrala.ag3nts.org/people"
const placesURL = "https://centrala.ag3nts.org/places"

type QueryRequest struct {
	APIKey string `json:"apikey"`
	Query  string `json:"query"`
}

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// "Kraków Łódź" -> ["KRAKOW", "LODZ"]
func sanitize(s string) []string {
	s = strings.ToUpper(s)
	s = strings.Map(func(r rune) rune {
		switch r {
		case 'Ą':
			return 'A'
		case 'Ę':
			return 'E'
		case 'Ć':
			return 'C'
		case 'Ń':
			return 'N'
		case 'Ś':
			return 'S'
		case 'Ł':
			return 'L'
		case 'Ź':
			return 'Z'
		case 'Ż':
			return 'Z'
		case 'Ó':
			return 'O'
		default:
			return r
		}
	}, s)
	return strings.Split(s, " ")
}

func askAPI(url string, question string) []string {
	key := api.ApiKey()
	query := QueryRequest{
		APIKey: key,
		Query:  question,
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("failed to marshal query: %+v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to send POST request: %+v", err)
	}
	defer resp.Body.Close()

	var apiResponse ApiResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %+v", err)
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("failed to unmarshal response body: %+v", err)
	}
	if apiResponse.Code != 0 {
		log.Printf("%+v. Q: %s\n", apiResponse, question)
		return nil
	}
	r := sanitize(apiResponse.Message)
	return r
}

func getUniqueItems(items []string, known map[string]bool) []string {
	result := []string{}
	for _, item := range items {
		if _, found := known[item]; !found {
			result = append(result, item)
			known[item] = true
		}
	}
	return result
}

func extractPeopleAndPlaces(filePath string) ([]string, []string) {
	const system = `Zajmujesz się wyodrębnianiem słów imion i miejscowości z tekstu`

	// if file with extracted keywords already exists - just read it and move on to the next file
	extractedFilePath := filePath + ".extracted"
	if _, err := os.Stat(extractedFilePath); err == nil {
		content, err := os.ReadFile(extractedFilePath)
		peoplePlaces := strings.Split(string(content), "\n")
		if err != nil {
			log.Printf("failed to read people places file %s: %+v\n", extractedFilePath, err)
		} else {
			return sanitize(peoplePlaces[0]), sanitize(peoplePlaces[1])
		}
	}

	// read the input file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file %s: %+v\n", filePath, err)
	}

	// extract the people and places from content
	log.Print("Extracting people and places for", filePath)
	user, err := api.BuildPrompt("prompt.txt", "{{INPUT}}", string(content))
	if err != nil {
		log.Fatalf("failed to build prompt: %+v\n", err)
	}
	resultString, err := openai.CompletionCheap(user, system, nil)
	if err != nil {
		log.Fatalf("openai returend error: %+v", err)
	}

	peoplePlaces := strings.Split(resultString, "\n")
	if len(peoplePlaces) != 2 {
		log.Fatal("invalid people-places:", resultString)
	}

	// save person name: keywords to file
	err = os.WriteFile(extractedFilePath, []byte(resultString), os.ModePerm)
	if err != nil {
		log.Printf("failed to save keywords to file %s: %+v\n", extractedFilePath, err)
	}
	people := strings.ReplaceAll(peoplePlaces[0], " ", "")
	places := strings.ReplaceAll(peoplePlaces[1], " ", "")
	return strings.Split(people, ","), strings.Split(places, ",")
}

func main() {
	err := api.DownloadIfDoesntExistYet(barbaraURL, "downloads/barbara.txt")
	api.FatalOnError(err)

	// collect people and places we managed to discover
	knownPeople := map[string]bool{}
	knownPlaces := map[string]bool{}

	// get the initial people and places to process
	peopleToProcess, placesToProcess := extractPeopleAndPlaces("downloads/barbara.txt")

	// take turns and ask for places related to person and for people related to place until there is anything not processed yet
	for len(peopleToProcess) > 0 || len(placesToProcess) > 0 {
		// get new places from people
		for _, p := range peopleToProcess {
			// 1. ask for places
			places := askAPI(peopleURL, p)

			// 2. filter out known places and add new places to known place
			newPlaces := getUniqueItems(places, knownPlaces)

			// 3. add unique places to new places for processing
			placesToProcess = append(placesToProcess, newPlaces...)
		}
		peopleToProcess = nil // all peopleToProcess processed

		// get new people from places
		for _, p := range placesToProcess {
			// 1. ask for people
			people := askAPI(placesURL, p)

			// 2. filter out known people and add new people to known people
			newPeople := getUniqueItems(people, knownPeople)

			// 3. add unique people to new people for processing
			peopleToProcess = append(peopleToProcess, newPeople...)
		}
		placesToProcess = nil // all placesToProcess processed
	}

	// Verify answer
	for p := range knownPlaces {
		fmt.Print("checking ", p, "...")
		answer := p
		result, err := api.VerifyTaskAnswer("loop", answer, api.VerificationURL)
		if err != nil {
			fmt.Println("Answer verification failed:", err)
		} else {
			fmt.Println(result)
		}
	}
}
