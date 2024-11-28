package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = "Odpowiedz na pytanie w oparciu o załączoną treść strony HTML. Jeśli stronia nie zawiera odpowiedzi na pytanie, napisz [BRAK]. HTML:\n"

func main() {
	openai.Debug = true

	// get tasks
	url := "https://centrala.ag3nts.org/data/" + api.ApiKey() + "/softo.json"
	taskStr, err := api.FetchData(url)
	if err != nil {
		log.Fatal(err)
	}
	tasks := map[string]string{}
	err = json.Unmarshal([]byte(taskStr), &tasks)
	if err != nil {
		log.Fatal("Error deserializing JSON:", err)
	}

	// prepare answers collection
	answers := map[string]string{}

	// run crawler
	c := colly.NewCollector()
	c.MaxDepth = 10
	c.AllowedDomains = []string{"softo.ag3nts.org"}
	c.AllowURLRevisit = false
	c.CacheDir = "./cache"
	c.IgnoreRobotsTxt = true
	visitCount := 0
	totalSize := 0

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	// Check the page contents
	c.OnResponse(func(r *colly.Response) {
		html := string(r.Body)

		// avoid traps
		if strings.Contains(html, "ANTY BOT") {
			return
		}

		// collect&print diagnostics
		visitCount++
		size := len(r.Body)
		totalSize += size
		fmt.Printf("%2d: size=%d, total=%d, depth=%d, url=%s\n", visitCount, size, totalSize, r.Request.Depth, r.Request.URL)

		// expenses security gate
		if totalSize > 100000 {
			panic("money saving emergency STOP")
		}

		// try to extract answers from the page
		for id, question := range tasks {
			// skip if already answered this question
			if _, found := answers[id]; found {
				continue
			}

			// process question
			fmt.Println("Question:", question)
			rsp, err := openai.CompletionCheap(question+" Odpowiedz tak krótko jak to tylko mozliwe!", system+html, nil)
			if err != nil {
				fmt.Printf("openai failed: %+v\n", err.Error())
				continue
			}
			if rsp != "[BRAK]" {
				answers[id] = rsp
				fmt.Println("Num found answers:", len(answers))
			} else {
				fmt.Println("Response:", rsp)
			}
			if len(answers) == len(tasks) {
				c.AllowedDomains = []string{""} // finish web crawling
			}
		}

		// pause after every page to see what's happening
		fmt.Println("Press 'Enter' to continue...")
		fmt.Scanln()
	})
	c.Visit("https://softo.ag3nts.org")

	// verify the answer
	fmt.Println("answer:\n", answers)
	result, err := api.VerifyTaskAnswer("softo", answers, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
