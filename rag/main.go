package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/mateuszmidor/AiStudy/rag/llm"
	"github.com/mateuszmidor/AiStudy/rag/vecdb"
)

// knowledge is a list of strings that are used to train the vector db
var knowledge = []string{
	"Python is kind of snake",
	"C++ is programming language that produces fast programs",
	"Python is lame programming language",
	"Rust is programming language that produces robust programs",
	"Monty Python is a comedy show",
}

// questions is a list of questions that are used to test the RAG
var questions = []string{
	"Which programming language is fast?",
	"Which programming language is robust?",
	"What is Python?",
	"Who is lame?",
	"What comedy shows do you know?",
	"What programming languages do you know?",
	"What animals do you know?",
}

func main() {
	// fill vector db with knowledge
	slog.Info("feeding the retriever, can take a dozen seconds...")
	vecdb.FeedDB(knowledge)

	// ask questions
	for _, question := range questions {
		// retrieve information relevant to the question from vector db
		slog.Info("retrieving information regarding: " + question)
		rsp := vecdb.AskDB(question, 3)
		slog.Info("retrieved", "results", rsp)

		// create prompt that includes the retrieved information for ollama
		prompt := makePrompt(question, rsp)
		slog.Info("prompt: \n" + prompt) // multiline

		// generate response
		slog.Info("sending prompt to ollama...")
		response := llm.OllamaGenerateCompletion(prompt)
		slog.Info("response: " + response)
		fmt.Println()
		fmt.Println()
	}
}

// makePrompt creates a prompt for the ollama based on the question and the information pieces
func makePrompt(question string, informationPieces []vecdb.SearchResult) string {
	instruction := "Instruction: Based only on the provided information, answer the question in one short sentence."
	information := collectInformationPieces(informationPieces)
	question = "Question: " + question
	return instruction + "\n" + information + "\n" + question
}

// collectInformationPieces collects information pieces from the search results,
// checks if they are useful and returns them as a single string
func collectInformationPieces(informationItems []vecdb.SearchResult) string {
	var info []string
	for _, r := range informationItems {
		if isUsefulInformation(r) {
			info = append(info, "Information: "+r.Text)
		}
	}
	return strings.Join(info, "\n")
}

// isUsefulInformation checks if the information piece is usefulbased on the score and returns true if it is
func isUsefulInformation(r vecdb.SearchResult) bool {
	return r.Score > 0.25 // arbitrary threshold but should work fine for now
}
