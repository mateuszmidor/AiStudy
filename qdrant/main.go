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

// SearchResponse represents search response structure.
// Example response:
// {"result":[{"id":"9b31733d-aa7a-07e9-71a1-dd8110a83374","version":2,"score":0.7733528,"payload":{"text":"C++ is programming language that produces fast programs"}}],"status":"ok","time":0.001875241}
type SearchResponse struct {
	Result []struct {
		ID      string  `json:"id"`
		Version int     `json:"version"`
		Score   float64 `json:"score"`
		Payload struct {
			Text string `json:"text"`
		} `json:"payload"`
	} `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

const dbBaseURL = "http://localhost:6333/collections/"
const collectionName = "knowledge"

// addCollection creates new collection of entries in vector database
func addCollection(dimensions int) error {
	slog.Info("add collection", slog.String("name", collectionName), slog.Int("dimensions", dimensions))

	// Prepare database config
	config := CollectionConfig{
		Vectors: VectorConfig{
			Size:     dimensions,
			Distance: "Cosine",
		},
	}

	// Send request
	url := dbBaseURL + collectionName
	_, err := request(url, "PUT", config)

	// Return result
	return err
}

// addPoint adds new entry to collection
func addPoint(vector []float64, text string) error {
	slog.Info("add point", slog.String("text", text))

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
func search(vector []float64, text string) (*SearchResponse, error) {
	slog.Info("search", slog.String("text", text))

	// Prepare search query
	query := SearchQuery{Vector: vector, Limit: 1, WithPayload: true}

	// Send request
	url := dbBaseURL + collectionName + "/points/search"
	rspString, err := request(url, "POST", query)

	if err != nil {
		return nil, err
	}

	// Deserialize response
	var rsp *SearchResponse
	err = json.Unmarshal([]byte(rspString), &rsp)
	if err != nil {
		return nil, err
	}

	// Return result
	return rsp, nil
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
	// determine vector size for collection; depends on pre-trained model used for embeding
	dimensions := len(embed("Check embeding dimensions"))

	// create the collection in vector database
	panicOnError(addCollection(dimensions))

	// store embeddings in collection
	for _, k := range knowledge {
		embedding := embed(k)
		panicOnError(addPoint(embedding, k))
	}

	// ask database the questions
	for _, q := range questions {
		embedding := embed(q)
		response, err := search(embedding, q)
		panicOnError(err)
		for _, r := range response.Result {
			slog.Info("result", slog.String("payload", r.Payload.Text), slog.Float64("score", r.Score))
		}
	}
}
