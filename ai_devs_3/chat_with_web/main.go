package main

import (
	"bufio"
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

func articleToMarkDown(url string) string {
	// extract filename from url
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1] + ".md"

	// check if file exists in "downloads" directory
	filePath := "downloads/" + filename
	if _, err := os.Stat(filePath); err == nil {
		// file exists, return its contents
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("failed to read file %s: %+v", filePath, err)
		}
		return string(content)
	}

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

	if _, err := os.Stat("downloads"); os.IsNotExist(err) {
		err := os.Mkdir("downloads", os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create downloads directory: %+v", err)
		}
	}
	// save the MarkDown article for optional investigation
	err = os.WriteFile("downloads/"+filename, []byte(declutteredArticleMD), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to file article.md: %+v", err)
	}
	log.Println("saved page as:", filename)

	return declutteredArticleMD
}

func getQuestion() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(line)
}

func main() {
	url := getURLFromCmdLineArgs()
	declutteredArticleMD := articleToMarkDown(url)

	// question-answer loop
	for {
		fmt.Print("> ")
		question := getQuestion()
		system := "When answering the question, use only the information you can find in the Input Document. ONLY that information is allowed!"
		user := "Question: " + question + "\nInput Document:\n" + declutteredArticleMD
		openai.Debug = true
		answer, err := openai.CompletionCheap(user, system, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(answer)
		fmt.Println()
	}
}

func getURLFromCmdLineArgs() string {
	if len(os.Args) < 2 {
		log.Fatal("you need to provide input URL")
	}
	url := os.Args[1]
	return url
}
