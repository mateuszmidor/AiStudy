package main

import (
	// "bytes"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
)

// CollectionConfig represents the complete collection (Database) configuration.
type CollectionConfig struct {
	Vectors VectorConfig `json:"vectors"`
}

// VectorConfig represents the configuration for vectors.
type VectorConfig struct {
	Size     int    `json:"size"`     // how many dimensions
	Distance string `json:"distance"` // distance func ["Dot", "Cosine"]
}

// Point is a single entry in collection
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"` // optional
}

// Points represents multiple entries in collection
type Points struct {
	Points []Point `json:"points"`
}

// SearchQuery represents the search query structure.
type SearchQuery struct {
	Vector      []float64 `json:"vector"`       // search input
	Limit       int       `json:"limit"`        // how many entries to return?
	WithPayload bool      `json:"with_payload"` // should return payload text?
	WithVectors bool      `json:"with_vectors"`
}

const dbBaseURL = "http://localhost:6333/collections/"

// addCollection creates new collection of entries in vector database
func addCollection(name string, dimensions int) error {
	slog.Info("add collection", slog.String("name", name), slog.Int("dimensions", dimensions))

	// Prepare database config
	config := CollectionConfig{
		Vectors: VectorConfig{
			Size:     dimensions,
			Distance: "Cosine",
		},
	}

	// Send request
	url := dbBaseURL + name
	_, err := request(url, "PUT", config)

	// Return result
	return err
}

// addPoint adds new entry to collection
func addPoint(collectionName string, vector []float64, text string) error {
	slog.Info("add point", slog.String("collection_name", collectionName), slog.String("text", text))

	// Prepare payload
	payload := map[string]interface{}{
		"text": text,
	}
	point := Point{
		ID:      generateMD5HashString(text),
		Vector:  vector,
		Payload: payload,
	}
	points := Points{
		Points: []Point{point},
	}

	// Send request
	url := dbBaseURL + collectionName + "/points"
	_, err := request(url, "PUT", points)

	// Return result
	return err
}

// search looks up database entries similar to provided vector
func search(collectionName string, vector []float64, text string) (string, error) {
	slog.Info("search", slog.String("collection_name", collectionName), slog.String("text", text))

	// Prepare search query
	query := SearchQuery{Vector: vector, Limit: 1, WithPayload: true}

	// Send request
	url := dbBaseURL + collectionName + "/points/search"
	rsp, err := request(url, "POST", query)

	// Return result
	return rsp, err
}

// request is a helper func that sends http request to url with provided data, and returns the response as text
func request(url, method string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to perform request to %q: %s, Response: %s", url, resp.Status, bodyString)
	}

	return bodyString, nil
}

// generateMD5HashString generates an MD5 hash string from the provided text.
func generateMD5HashString(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

// panicOnError checks if the provided error is not nil and panics with that error.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// embed executes a Python script to generate an embedding for the given input string and returns it as a slice of float64 values.
func embed(input string) []float64 {
	cmd := exec.Command("python", "./embedding-localhost/main.py", input)
	output, err := cmd.Output()
	panicOnError(err)
	var embedding []float64
	panicOnError(json.Unmarshal(output, &embedding))
	return embedding
}

var knowledge = []string{
	"Python is kind of snake",
	"Python is lame programming language",
	"C++ is programming language that produces fast programs",
	"Rust is programming language that produces robust programs",
}
var questions = []string{
	"Which programming language is fast?",
	"Which programming language is robust?",
	"What is Python?",
	"Who is lame?",
}

func main() {
	// determine vector size for collection; depends on pre-trained model used for embeding
	dimensions := len(embed("Check embeding dimensions"))

	// create the collection in vector database
	panicOnError(addCollection("knowledge", dimensions))

	// store embeddings in collection
	for _, s := range knowledge {
		embedding := embed(s)
		panicOnError(addPoint("knowledge", embedding, s))
	}

	// ask database the questions
	for _, q := range questions {
		embedding := embed(q)
		result, err := search("knowledge", embedding, q)
		panicOnError(err)
		fmt.Println(result)
	}
}
