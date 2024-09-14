package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/mateuszmidor/AiStudy/rag/llm"
	"github.com/mateuszmidor/AiStudy/rag/vecdb"
)

var knowledge = []string{
	"Python is kind of snake",
	"C++ is programming language that produces fast programs",
	"Python is lame programming language",
	"Rust is programming language that produces robust programs",
	"Monty Python is a comedy show",
}

var questions = []string{
	"Which programming language is fast?",
	"Which programming language is robust?",
	"What is Python?",
	"Who is lame?",
	"Do you know any comedy shows?",
}

func main() {
	// fill vector db with knowledge
	vecdb.CreateDB(knowledge)

	// ask questions
	for _, q := range questions {
		// retrieve information relevant to the question from vector db
		rsp := vecdb.AskDB(q, 3)
		slog.Info("search", "results", rsp)

		// create prompt that includes the retrieved information for ollama
		prompt := makePrompt(q, rsp)

		// generate response
		response := llm.OllamaGenerateCompletion(prompt)
		slog.Info("Response:\n" + response)
		fmt.Println()
		fmt.Println()
	}
}

func makePrompt(question string, informationItems []vecdb.SearchResult) string {
	prompt := "Based only on the provided information, briefly answer the question.\nQuestion: " + question + "\n"

	var info []string
	for _, r := range informationItems {
		if isUsefulInformation(r) {
			info = append(info, "Information: "+r.Text)
		}
	}

	return prompt + strings.Join(info, "\n")
}

func isUsefulInformation(r vecdb.SearchResult) bool {
	return r.Score > 0.25 // arbitrary threshold but should work fine for now
}
