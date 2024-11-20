package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func declutteredArticleMD(inputMD string) string {
	// Remove links from the markdown content
	declutteredMD := inputMD
	declutteredMD = removeLinks(declutteredMD)
	return declutteredMD
}

// Helper function to remove links from markdown
func removeLinks(md string) string {
	lines := strings.Split(md, "\n")
	var filteredLines []string
	for _, line := range lines {
		// Use a regular expression to find and remove markdown links
		re := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
		lineWithoutLinks := re.ReplaceAllString(line, "$1")
		filteredLines = append(filteredLines, lineWithoutLinks)
	}
	return strings.Join(filteredLines, "\n")
}

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

	declutteredArticleMD := declutteredArticleMD(articleMD)

	// save the MarkDown article for optional investigation
	err = os.WriteFile("article.md", []byte(declutteredArticleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to file article.md: %+v", err)
	}

	// answer the question
	system := "When answering the question, use only the information you can find in the Input Document. ONLY that information is allowed!"
	user := "Question: " + question + "\nInput Document:\n" + declutteredArticleMD
	openai.Debug = true
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
