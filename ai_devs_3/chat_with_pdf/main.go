package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

func main() {
	pdfFilename, question := getFilenameAndQuestionFromCmdLineArgs()
	txtFilename := convertPdfToTxt(pdfFilename)
	txtContent := readFile(txtFilename)

	system := "When answering the question, use only the information you can find in the Input Document. ONLY that information is allowed!"
	user := "Question: " + question + "\nInput Document:\n" + txtContent
	answer, err := openai.CompletionCheap(user, system, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)
}

func readFile(txtFilename string) string {
	contentBytes, err := os.ReadFile(txtFilename)
	if err != nil {
		log.Fatalf("failed to read text file: %v", err)
	}
	content := string(contentBytes)
	return content
}

func convertPdfToTxt(pdfFilename string) string {
	txtFilename := pdfFilename + ".txt"
	if _, err := os.Stat(txtFilename); err == nil {
		return txtFilename
	}
	cmd := fmt.Sprintf("pdftotext -layout %s %s", pdfFilename, txtFilename)
	err := runCommand(cmd)
	if err != nil {
		log.Fatalf("failed to convert PDF to text: %v", err)
	}
	return txtFilename
}

func getFilenameAndQuestionFromCmdLineArgs() (string, string) {
	if len(os.Args) < 3 {
		log.Fatal("you need to provide input PDF and a question")
	}
	pdfFilename := os.Args[1]
	question := os.Args[2]
	return pdfFilename, question
}

// runCommand executes a shell command and returns an error if it fails
func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
