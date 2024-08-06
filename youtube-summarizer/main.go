package main

import (
	"fmt"
	"strings"
)

func main() {
	videoURL := "https://www.youtube.com/watch?v=XOKtkMVW_lY" // must be a video with captions
	captions, err := getCaptions(videoURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	joinedCaptions := strings.Join(captions, "\n")
	prompt := fmt.Sprintf("Summarize the following text in bullet point format, you MUST respond in the same language as the source text.\nSource text:\n###\n%s\n###", joinedCaptions)
	completion := ollamaGenerateCompletion(prompt)
	fmt.Println(completion)
}
