package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

// place the input document in system prompt at the very beginning, so OpenID prompt caching is used
const system = "When answering the question, use only the information you can find in the Input Document. ONLY that information is allowed! Input Document:\n"

func main() {
	url := getURLFromCmdLineArgs()
	declutteredWebPageMD := webPageToMarkDown(url)
	chat := prepareChat(system + declutteredWebPageMD)

	// question-answer loop
	for {
		fmt.Print("> ")
		question := getQuestion()
		user := question
		answer, err := chat.User(user, nil, nil, "text", 0)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(answer.Choices[0].Message.Content)
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

func declutteredArticleMD(inputMD string) string {
	return removeLinks(inputMD)
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

func webPageToMarkDown(pageURL string) string {
	filename := extractFilename(pageURL)
	filepath := os.TempDir() + filename
	if content, err := readFileIfExists(filepath); err == nil {
		return content
	}

	articleHTML := fetchArticleHTML(pageURL)
	articleMD := convertHTMLToMarkdown(articleHTML)
	declutteredMD := declutteredArticleMD(articleMD)

	saveMarkdownToFile(declutteredMD, filepath)

	return declutteredMD
}

func extractFilename(pageURL string) string {
	pageURL = strings.TrimSuffix(pageURL, "/")
	parts := strings.Split(pageURL, "/")
	filename := parts[len(parts)-1] + ".md"
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		log.Fatalf("failed to url-decode filename: %+v", err)
	}
	return decodedFilename
}

func readFileIfExists(filePath string) (string, error) {
	if _, err := os.Stat(filePath); err == nil {
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("failed to read file %s: %+v", filePath, err)
		}
		log.Println("read page from:", filePath)
		return string(content), nil
	}
	return "", fmt.Errorf("file does not exist")
}

func fetchArticleHTML(pageURL string) string {
	articleHTML, err := api.FetchData(pageURL)
	if err != nil {
		log.Fatalf("Error fetching questions: %+v", err)
	}
	return articleHTML
}

func convertHTMLToMarkdown(articleHTML string) string {
	articleMD, err := htmltomarkdown.ConvertString(articleHTML)
	if err != nil {
		log.Fatal(err)
	}
	return articleMD
}

func saveMarkdownToFile(markdownContent, filename string) {
	err := os.WriteFile(filename, []byte(markdownContent), os.ModePerm)
	if err != nil {
		log.Fatalf("failed to save article to %s: %+v", filename, err)
	}
	log.Println("saved page as:", filename)
}

func getQuestion() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(line)
}

func prepareChat(system string) *openai.ChatWithMemory {
	openai.Debug = true
	chat, err := openai.NewChatWithMemory(system, "gpt-4o-mini", 15000, false)
	if err != nil {
		log.Fatal(err)
	}
	return chat
}
