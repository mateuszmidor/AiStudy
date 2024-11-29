package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const system = `
Jesteś pionkiem na szachownicy o rozmiarach 4x4, stoisz na polu 1,1. 
Twoim zadaniem jest odpowiedzieć na jakim polu stoisz po wykonaniu polecenia od uzytkownika.
Odpowiadasz w formacie JSON z polami: _thinking, x, y. _thinking musi być pierwszym polem.
Przykład:
	Polecenie: przesuwasz się o jeden w prawo
	Odpowiedź: 
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się o 1 w prawo, czyli na pozycję 2,1.",
		"x": "2",
		"y": "1"
	}

	Polecenie: przesuwasz się o dwa w prawo
	Odpowiedź: 
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się o 2 w prawo, czyli na pozycję 3,1.",
		"x": "3",
		"y": "1"
	}

	Polecenie: przesuwasz się o jeden w dół
	Odpowiedź: 
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się o 1 w dół, czyli na pozycję 1,2.",
		"x": "1",
		"y": "2"
	}

	Polecenie: przesuwasz się o dwa w dół
	Odpowiedź: 
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się o 2 w dół, czyli na pozycję 1,3.",
		"x": "1",
		"y": "3"
	}

	Polecenie: przesuwasz się na maksa w prawo i potem o jeden w dół
	Odpowiedź: 
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się na maksa w prawo, czyli na 4,1. Potem jeden w dół, czyli na 4,2",
		"x": "4",
		"y": "2"
	}

	Polecenie: przesuwasz się na maksa w dół i potem o jeden w prawo
	Odpowiedź:
	{
		"_thinkinig" : "startuję z pozycji 1,1. przesuwam się o na maksa w dół, czyli na 1,4. Potem o jeden w prawo, czyli na 2,4",
		"x": "2",
		"y": "4"
	}
`

type llmresponse struct {
	Thinking string `json:"_thinking"`
	X        string `json:"x"`
	Y        string `json:"y"`
}
type webhookRequest struct {
	Instruction string `json:"instruction"`
}
type webhookResponse struct {
	Description string `json:"description"`
}

// ........ y,x
var theMap = [][]string{
	{"start", "łąka", "drzewo", "dom"},
	{"łąka", "wiatrak", "łąka", "łąka"},
	{"łąka", "łąka", "skały", "drzewo"},
	{"skały", "skały", "auto", "jaskinia"},
}

func main() {
	// run webhook server
	go runWebHookServer()

	// print the map
	for _, row := range theMap {
		for _, cell := range row {
			fmt.Printf("%-10s", cell) // Adjust the width as needed
		}
		fmt.Println()
	}
	fmt.Println()

	// run the conversation in repeat until success fashion:)
	externalURL := os.Args[1]
	time.Sleep(time.Second)
	for {
		fmt.Println("Press 'Enter' to continue...")
		fmt.Scanln()
		rsp, err := api.VerifyTaskAnswer("webhook", externalURL, api.VerificationURL)
		if err != nil {
			log.Println("Error:", err)
		} else {
			fmt.Println(rsp)
		}
		time.Sleep(time.Second)
	}
}

func runWebHookServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// decode the request
		decoder := json.NewDecoder(r.Body)
		var req webhookRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Println("Invalid JSON.", err.Error())
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		log.Println("received:", req.Instruction)

		// find out the final position after traveling as instructed
		gptrsp, err := openai.CompletionExpert(req.Instruction, system, nil, nil, "gpt-4o-mini", "json_object", 100, 0)
		if err != nil {
			log.Printf("openai failed: %+v\n", err)
			return
		}
		rsp := gptrsp.Choices[0].Message.Content
		log.Println("openai responded with:", rsp)

		var answ llmresponse
		err = json.Unmarshal([]byte(rsp), &answ)
		if err != nil {
			log.Println("Error deserializing response:", err.Error())
			http.Error(w, "Error deserializing response", http.StatusInternalServerError)
			return
		}
		x, _ := strconv.Atoi(answ.X)
		y, _ := strconv.Atoi(answ.Y)
		x--
		y--

		// respond with object from the map
		log.Println("responding with:", theMap[y][x])
		hookRsp := webhookResponse{Description: theMap[y][x]}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(hookRsp)
		if err != nil {
			log.Println("Error encoding response to JSON.", err.Error())
			http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
			return
		}
		time.Sleep(time.Second)
	})

	fmt.Println("Starting server on port 33000...")
	err := http.ListenAndServe(":33000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
