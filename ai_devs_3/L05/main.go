package main

import (
	"fmt"
	"log"
	"os"

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
Remember: IT IS CRUTIAL that you don't change anything other in the text than the sensitive information and only return the redacted text. 
Text to redact:
`

func main() {
	// fetch the task input
	apikey := os.Getenv("AIDEVS3_API_KEY")
	taskUrl := "https://centrala.ag3nts.org/data/" + apikey + "/cenzura.txt"
	personalData, err := api.FetchTask(taskUrl)
	if err != nil {
		log.Fatalf("Error fetching task: %+v", err)
	}

	// use local LLM to censor personal information
	prompt := promptTemplate + personalData
	censoredData, err := ollama.Completion(prompt, "llama3")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Uncensored: %q\n", personalData)
	fmt.Printf("Censored:   %q\n", censoredData)

	// verify the answer
	result, err := api.VerifyTaskAnswer("CENZURA", censoredData, verifyURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
