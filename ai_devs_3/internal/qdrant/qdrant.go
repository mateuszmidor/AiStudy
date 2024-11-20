package qdrant

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// CollectionConfig represents the complete collection (Database) configuration.
type CollectionConfig struct {
	Vectors VectorConfig `json:"vectors"`
}

// VectorConfig represents the configuration for vectors.
type VectorConfig struct {
	Size     int    `json:"size"`     // how many dimensions
	Distance string `json:"distance"` // distance func ["Cosine", "Dot", "Euclidean", "Manhattan"]
}

// Point is a single entry in collection
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"` // optional, you can store the original text here and some meta info/tags describing the content
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
		ID      string                 `json:"id"`
		Version int                    `json:"version"`
		Score   float64                `json:"score"`
		Payload map[string]interface{} `json:"payload"`
	} `json:"result"`
	Status string  `json:"status"`
	Time   float64 `json:"time"`
}

type Question struct {
	Vector []float64
	// filters - to implement
}

type SearchResult struct {
	Score   float64                // range 0-1
	Payload map[string]interface{} // payload of the found vector db entry
}

const dbBaseURL = "http://localhost:6333/collections/"
const collectionName = "knowledge"

// FeedDB creates a new collection in the vector database and stores the provided knowledge in form of embeddings
func FeedDB(points []Point, dimensions int) error {
	// create the collection in vector database
	err := addCollection(dimensions)
	if err != nil {
		return err
	}

	// store embeddings in collection
	for _, p := range points {
		err := addPoint(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// AskDB retrieves information from the vector database based on the provided question, it returns a maximum of maxAnswers
func AskDB(question Question, maxAnswers int) (result []SearchResult, err error) {
	response, err := search(question, maxAnswers)
	if err != nil {
		return nil, err
	}

	for _, r := range response.Result {
		result = append(result, SearchResult{Score: r.Score, Payload: r.Payload})
	}

	return result, nil
}

// addCollection creates new collection of entries in vector database
func addCollection(dimensions int) error {

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
func addPoint(p Point) error {
	// Prepare payload
	point := Point{
		ID:      generateMD5HashString(p.Payload["text"].(string)),
		Vector:  p.Vector,
		Payload: p.Payload,
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
func search(question Question, maxAnswers int) (*SearchResponse, error) {
	// Prepare search query
	query := SearchQuery{Vector: question.Vector, Limit: maxAnswers, WithPayload: true}

	// Send request
	url := dbBaseURL + collectionName + "/points/search"
	rspString, err := request(url, "POST", query)

	if err != nil {
		return nil, errors.Wrap(err, "qdrant request failed")
	}

	// Deserialize response
	var rsp *SearchResponse
	err = json.Unmarshal([]byte(rspString), &rsp)
	if err != nil {
		return nil, errors.Wrap(err, "response unmarshall failed")
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
