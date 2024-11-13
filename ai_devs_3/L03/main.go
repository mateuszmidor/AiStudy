package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const verifyURL = "https://centrala.ag3nts.org/report"

type TestData struct {
	Question string `json:"question"`
	Answer   int    `json:"answer"`
	Test     *Test  `json:"test,omitempty"`
}

type Test struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Data struct {
	Apikey      string     `json:"apikey"`
	Description string     `json:"description"`
	Copyright   string     `json:"copyright"`
	TestData    []TestData `json:"test-data"`
}

func main() {
	// Read the file from disk
	file, err := os.Open("json.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Verify all items
	for i, item := range data.TestData {
		// Correct calculation errors
		parts := strings.Split(item.Question, " + ")
		if len(parts) == 2 {
			num1, err1 := strconv.Atoi(parts[0])
			num2, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				correctAnswer := num1 + num2
				if item.Answer != correctAnswer {
					fmt.Printf("Correcting answer for question '%s': %d -> %d\n", item.Question, item.Answer, correctAnswer)
					data.TestData[i].Answer = correctAnswer
				}
			}
		}

		// Handle test questions
		if item.Test != nil && item.Test.A == "???" {
			answer, err := openai.CompletionCheap(item.Test.Q, "", nil)
			if err != nil {
				fmt.Printf("Error getting answer for question '%s': %v\n", item.Test.Q, err)
			} else {
				fmt.Printf("Answering question '%s' -> %s\n", item.Test.Q, answer)
				data.TestData[i].Test.A = answer
			}
		}
	}

	// Output the corrected data
	data.Apikey = os.Getenv("AIDEVS3_API_KEY")
	result, err := api.VerifyTaskAnswer("JSON", data, verifyURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
