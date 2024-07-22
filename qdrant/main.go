package main

import (
	// "bytes"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FullConfig represents the complete configuration including the "vectors" object.
type FullConfig struct {
	Vectors VectorConfig `json:"vectors"`
}

// VectorConfig represents the configuration for vectors.
type VectorConfig struct {
	Size     int    `json:"size"`
	Distance string `json:"distance"`
}

// Points represents the structure of each item in the "points" array.
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

// Points represents the structure of the entire JSON object.
type Points struct {
	Points []Point `json:"points"`
}

// SearchQuery represents the search query structure.
type SearchQuery struct {
	Vector      []float64 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
	WithVectors bool      `json:"with_vectors"`
}

func addCollection(name string) error {
	// Example usage
	config := FullConfig{
		Vectors: VectorConfig{
			Size:     4,
			Distance: "Cosine",
		},
	}

	// Serialize the config to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// Define the URL for the collection
	url := "http://localhost:6333/collections/" + name

	// Create a new request using http
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add collection: %s", resp.Status)
	}

	return nil

}

func addPoint(collectionName string, vector []float64, text string) error {
	payload := map[string]interface{}{
		"text": text,
	}
	point := Point{
		ID:      generateMD5HashString(text),
		Vector:  vector,
		Payload: payload,
	}
	col := Points{
		Points: []Point{point},
	}

	jsonData, err := json.Marshal(col)
	if err != nil {
		return err
	}

	url := "http://localhost:6333/collections/" + collectionName + "/points"

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to add point: %s, Response: %s", resp.Status, bodyString)
	}

	return nil
}

func generateMD5HashString(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func search(collectionName string, vector []float64) error {
	q := SearchQuery{Vector: vector, Limit: 2, WithPayload: true}

	jsonData, err := json.Marshal(q)
	if err != nil {
		return err
	}

	url := "http://localhost:6333/collections/" + collectionName + "/points/search"

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("failed to perform search: %s, Response: %s", resp.Status, bodyString)
	}

	// Print the response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	return nil
}

// PanicOnError checks if the provided error is not nil and panics with that error.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	PanicOnError(addCollection("dane"))
	PanicOnError(addPoint("dane", []float64{-1, 0, 0, 1}, "ala ma kota"))
	PanicOnError(addPoint("dane", []float64{-1, 0, 0, 0.9}, "ala ma psa"))
	time.Sleep(time.Millisecond * 250)
	PanicOnError(search("dane", []float64{-1, 0, 0, 1}))
}
