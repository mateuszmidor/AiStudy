package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const databaseApiURL = "https://centrala.ag3nts.org/apidb"

// return pairs of connected person names
const connectionsQuery = `
SELECT name1.username AS name1, name2.username AS name2
FROM users AS name1
JOIN connections AS c ON name1.id = c.user1_id
JOIN users AS name2 ON name2.id = c.user2_id
`

type DatabaseQuery struct {
	Task   string `json:"task"`
	APIKey string `json:"apikey"`
	Query  string `json:"query"`
}

type ConnectionsReply struct {
	Reply []Connection `json:"reply"`
}

type Connection struct {
	Name1 string `json:"name1"`
	Name2 string `json:"name2"`
}

func askAPI(question string) string {
	key := api.ApiKey()
	query := DatabaseQuery{
		Task:   "database",
		APIKey: key,
		Query:  question,
	}

	jsonData, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("failed to marshal query: %+v", err)
	}

	resp, err := http.Post(databaseApiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to send POST request: %+v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %+v", err)
	}

	result := string(body)
	// fmt.Println(result)
	return result
}

func insecureNeo4jConfig(c *neo4j.Config) {
	c.Encrypted = false
}

func main() {
	// grab connections from remote database
	connectionsJSON := askAPI(connectionsQuery)
	var connections ConnectionsReply
	err := json.Unmarshal([]byte(connectionsJSON), &connections)
	if err != nil {
		log.Fatalf("failed to unmarshal connectionsJSON: %+v", err)
	}

	// feed Neo4j with connections
	driver, err := neo4j.NewDriver("neo4j://localhost:7687", neo4j.NoAuth(), insecureNeo4jConfig)
	if err != nil {
		log.Fatalf("failed to create driver: %+v", err)
	}
	defer driver.Close()

	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer session.Close()

	// clear all existing nodes and edges, if any
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	if err != nil {
		log.Fatalf("failed to clear all nodes and edges: %+v", err)
	}

	// create people (nodes)
	uniquePeople := make(map[string]struct{})
	for _, connection := range connections.Reply {
		uniquePeople[connection.Name1] = struct{}{}
		uniquePeople[connection.Name2] = struct{}{}
	}
	for person := range uniquePeople {
		_, err = session.Run("CREATE (person:Person {name: $name})", map[string]interface{}{"name": person})
		if err != nil {
			log.Fatalf("failed to create Person node for %s: %+v", person, err)
		}
	}

	// create connection (edges)
	for _, c := range connections.Reply {
		query := fmt.Sprintf("MATCH (a:Person {name: '%s'}), (b:Person {name: '%s'}) CREATE (a)-[:Connection]->(b)", c.Name1, c.Name2)
		_, err = session.Run(query, nil)
		if err != nil {
			log.Fatalf("failed to create connection between Travolta and Smith: %+v", err)
		}
	}

	// read people names - just for testing
	allPeople := []string{}
	neo4jResult, err := session.Run("MATCH (p:Person) RETURN p.name", nil)
	if err != nil {
		log.Fatalf("failed to execute query: %+v", err)
	}
	for neo4jResult.Next() {
		name := neo4jResult.Record().GetByIndex(0).(string)
		allPeople = append(allPeople, name)
	}
	if err = neo4jResult.Err(); err != nil {
		log.Fatalf("result error: %+v", err)
	}
	fmt.Println("All people:", allPeople)

	// find shortest path from Rafał to Barbara
	shortestPathQuery := `
	MATCH (start:Person {name: 'Rafał'}), (end:Person {name: 'Barbara'}),
	path = shortestPath((start)-[*]-(end))
	RETURN path
`
	neo4jResult, err = session.Run(shortestPathQuery, nil)
	if err != nil {
		log.Fatalf("failed to find shortest path: %+v", err)
	}

	pathFound := neo4jResult.Next()
	if !pathFound {
		log.Fatal("failed to find path from Rafał to Barbara")
	}
	path := neo4jResult.Record().GetByIndex(0).(neo4j.Path)
	shortestPath := []string{}
	for _, p := range path.Nodes() {
		shortestPath = append(shortestPath, p.Props()["name"].(string))
	}
	fmt.Println("Shortest path:", shortestPath)

	answer := strings.Join(shortestPath, ", ")
	result, err := api.VerifyTaskAnswer("connections", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}
