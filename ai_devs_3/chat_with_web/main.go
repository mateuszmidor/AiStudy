package main

import (
	"fmt"
	"log"
	"os"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	url, question := getURLAndQuestionFromCmdLineArgs()

	// fetch the article
	articleHTML, err := api.FetchData(url)
	if err != nil {
		log.Fatalf("Error fetching questions: %+v", err)
	}

	// transform article HTML to MarkDown for LLM to understand it better
	articleMD, err := htmltomarkdown.ConvertString(articleHTML)
	if err != nil {
		log.Fatal(err)
	}

	// save the MarkDown article for optional investigation
	err = os.WriteFile("article.md", []byte(articleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to file article.md: %+v", err)
	}

	// answer the question
	system := "When answering the question, use only the information you can find in the Input Document. ONLY that information is allowed!"
	user := "Question: " + question + "\nInput Document:\n" + articleMD
	answer, err := openai.CompletionCheap(user, system, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)
}

func getURLAndQuestionFromCmdLineArgs() (string, string) {
	if len(os.Args) < 3 {
		log.Fatal("you need to provide input URL and a question")
	}
	url := os.Args[1]
	question := os.Args[2]
	return url, question
}
