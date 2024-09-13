package main

import (
	"log"
	"net/http"
	"os"
)

const address = "0.0.0.0:33000"

var counter int

func main() {
	http.HandleFunc("/", handler)
	log.Println("serving at", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	counter++
	log.Printf("%d - a request came!", counter)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(getResponseText()))
}

func getResponseText() string {
	data, err := os.ReadFile("response.txt")
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}
	return string(data)
}
