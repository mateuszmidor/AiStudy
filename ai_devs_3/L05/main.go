package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/ollama"
)

const verifyURL = "https://centrala.ag3nts.org/report"
const promptTemplate = `
<objective>
You are a precise sensitive-text redactor who replaces all occurence of personal sensitive information with the word CENZURA.
</objective>
<rules>
1. Replace every single personal sensitive information like first name, last name, city name, street name, street number, age, with the word CENZURA.
2. Keep all non-sensitive parts of the text intact. Keep punctuation intact. Only replace the sensitive parts.
3. Only return the redacted content WITHOUT ANY additional changes or extra text.
</rules>
<example_1>
input = Dane osoby podejrzanej: Paweł Zieliński. Zamieszkały w Warszawie na ulicy Pięknej 5. Ma 28 lat.
expected_output = Dane osoby podejrzanej: CENZURA. Zamieszkały w CENZURA na ulicy CENZURA. Ma CENZURA lat.
</example_1>
<example_2>
input = Informacje o podejrzanym: Marek Jankowski. Mieszka w Białymstoku na ulicy Lipowej 9. Wiek: 26 lat.
expected_output = Informacje o podejrzanym: CENZURA. Mieszka w CENZURA na ulicy CENZURA. Wiek: CENZURA lat.
</example_2>
<example_3>
input = Tożsamość podejrzanego: Michał Wiśniewski. Mieszka we Wrocławiu na ul. Słonecznej 20. Wiek: 30 lat.
expected_output = Tożsamość podejrzanego: CENZURA. Mieszka we CENZURA na ul. CENZURA. Wiek: CENZURA lat.
</example_3>
Remember: IT IS CRUTIAL that you don't change anything other in the text than the sensitive information. The life depends on it. I mean it!
Text to redact:
`

func main() {
	// fetch the personal data
	apikey := os.Getenv("AIDEVS3_API_KEY")
	httpResp, err := http.Get("https://centrala.ag3nts.org/data/" + apikey + "/cenzura.txt")
	if err != nil {
		log.Fatal("Error fetching text file:", err)
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal("Error reading http response body:", err)
	}
	dataStr := string(data)

	// use local LLM to hide personal information
	prompt := promptTemplate + dataStr
	ollamaResp, err := ollama.Completion(prompt, "llama3")
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("Uncensored: %q\n", dataStr)
	fmt.Printf("Censored:   %q\n", ollamaResp)

	// verify the answer
	result, err := api.VerifyTaskAnswer("CENZURA", ollamaResp, verifyURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
	Prompt string `json:"prompt"`
}

type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"` // prompt tokens
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"` // response tokens
	EvalDuration       int64     `json:"eval_duration"`
}

func ollamaGenerateCompletion(prompt string) string {
	fmt.Printf("Prompt:\n%s\n\n", prompt)

	// Initialize the payload
	payload := &OllamaRequest{
		Model:  "llama3",
		Stream: false,
		Prompt: prompt,
	}
	fmt.Println("num characters:", len(payload.Prompt))

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	// Specify the URL
	url := "http://localhost:11434/api/generate"

	// Create a new request using http.Post
	fmt.Println("Sending prompt to ollama...")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// Unmarshal the JSON response into an OllamaResponse struct
	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	fmt.Println("Received response from ollama:")
	fmt.Println("- input tokens:", ollamaResponse.PromptEvalCount)
	fmt.Println("- output tokens:", ollamaResponse.EvalCount)
	return ollamaResponse.Response
}
