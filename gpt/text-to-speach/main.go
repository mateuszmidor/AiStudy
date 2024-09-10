package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

const openAIURL = "https://api.openai.com/v1/audio/speech"

// GPTRequest represents the structure of the JSON for the text-to-speech request.
type GPTRequest struct {
	Model  string  `json:"model"`                     // REQUIRED [tts-1, tts-1-hd]
	Input  string  `json:"input"`                     // REQUIRED the actual text to be voiced out
	Voice  string  `json:"voice"`                     // REQUIRED [alloy, echo, fable, onyx, nova, shimmer]
	Format string  `json:"response_format,omitempty"` // [mp3, opus, aac, flac, wav, pcm], default: mp3
	Speed  float64 `json:"speed,omitempty"`           // [0.25-4.0], default: 1.0
}

// GPTErrorResponse represents the JSON error response structure from the GPT API.
type GPTErrorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// text length limit is 4096 characters
func textToSpeach(text string, outputMP3 string) {
	// get API KEY from env
	apiKey := os.Getenv("GPT_APIKEY")
	if apiKey == "" {
		panic("OpenAI API key is not set")
	}

	// send request
	req := prepareTextToSpeachRequest(text, apiKey)
	resp, err := http.DefaultClient.Do(req)
	panicOnError(err)
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	panicOnError(err)
	if resp.StatusCode != http.StatusOK {
		var errorResponse GPTErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		panicOnError(err)
		errorMsg := fmt.Sprintf("request failed [code:%d]: %s [type:%s]", resp.StatusCode, errorResponse.Error.Message, errorResponse.Error.Type)
		panic(errorMsg)
	}

	// write audio file
	err = os.WriteFile(outputMP3, body, fs.ModePerm)
	panicOnError(err)
}

// https://platform.openai.com/docs/api-reference/audio/createSpeech
func prepareTextToSpeachRequest(text string, apiKey string) *http.Request {
	// prepare request body
	reqBody := GPTRequest{
		Model: "tts-1",
		Input: text,
		Voice: "nova",
	}

	reqBytes, err := json.Marshal(reqBody)
	panicOnError(err)

	req, err := http.NewRequest("POST", openAIURL, bytes.NewReader(reqBytes))
	panicOnError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return req
}

// user panicOnError to reduce lines of code by 50%.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	const inputText = "Nosił wilk razy kilka, ponieśli i wilka."
	const oputputMP3 = "voiced.mp3"

	textToSpeach(inputText, oputputMP3)
	fmt.Println("Outputted audio file:", oputputMP3)
}
