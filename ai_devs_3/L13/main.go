package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const databaseApiURL = "https://centrala.ag3nts.org/apidb"

type DatabaseQuery struct {
	Task   string `json:"task"`
	APIKey string `json:"apikey"`
	Query  string `json:"query"`
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

// prepare functions for GPT to select from
func prepareTools() []openai.Tool {
	// prepare temp function parameters
	sqlStatementFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"sql_statement": {
				Type:        "string",
				Description: "sql statement to execute, e.g. 'select * from users'",
			},
		},
	}
	getTableSchemaFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"table_name": {
				Type:        "string",
				Description: "name of the sql database table to get the create statement for, e.g. 'users'",
			},
		},
	}

	// prepare functions
	sqlStatementFunc := openai.Function{
		Name:        "exec_sql_statement",
		Description: "this function executes given sql statement",
		Parameters:  &sqlStatementFuncParams,
	}
	getTableSchemaFunc := openai.Function{
		Name:        "get_sql_db_table_schema",
		Description: "this function takes sql database table name as a parameter and returns the exact sql statement that was used to create the table",
		Parameters:  &getTableSchemaFuncParams,
	}
	showTablesFunc := openai.Function{
		Name:        "show_sql_db_tables",
		Description: "this function returns a list of existing sql database tables",
		Parameters:  nil, // no parameters
	}

	// prepare tools from functions
	getTableSchema := openai.Tool{
		Type:     "function",
		Function: getTableSchemaFunc,
	}
	showTables := openai.Tool{
		Type:     "function",
		Function: showTablesFunc,
	}
	sqlStatement := openai.Tool{
		Type:     "function",
		Function: sqlStatementFunc,
	}

	// return tool set
	return []openai.Tool{sqlStatement, getTableSchema, showTables}
}

func getArgs(argsJSON string) map[string]string {
	var args map[string]string
	err := json.Unmarshal([]byte(argsJSON), &args)
	if err != nil {
		log.Fatalf("failed to unmarshal argsJSON: %+v", err)
	}
	return args
}

func main() {
	// askAPI("show tables")
	// askAPI("select * from datacenters")
	// askAPI("show create table datacenters")

	tools := prepareTools()
	chat, err := openai.NewChatWithMemory("Solve the task using available tools. Reply with [FINISHED] when the task is solved", "gpt-4o-mini", 1000)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	rsp, err := chat.User("które aktywne datacenter (DC_ID) są zarządzane przez pracowników, którzy są na urlopie (is_active=0)", nil, tools, "text", 0)
	// rsp, err := chat.User("create sql statement that returns active people from the relevant table. Use tools to check table structure", nil, tools, "text", 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	for {
		for _, choice := range rsp.Choices {
			if choice.Message.Content != "" {
				fmt.Println(choice.Message.Role, ":", choice.Message.Content)

			}

			for _, toolCal := range choice.Message.ToolCalls {
				name := toolCal.Function.Name
				args := getArgs(toolCal.Function.Arguments)
				callID := toolCal.ID
				fmt.Println("Tool", name, args, callID)
				result := ""
				switch name {
				case "exec_sql_statement":
					result = askAPI(args["sql_statement"])
				case "get_sql_db_table_schema":
					result = askAPI("show create table " + args["table_name"])
				case "show_sql_db_tables":
					result = askAPI("show tables")
				default:
					log.Fatal("unknown function: " + name)
				}
				rsp, err = chat.ToolResponse(result, callID)
				if err != nil {
					log.Fatalf("%+v", err)
				}
				fmt.Println(rsp.Choices[0].Message.Content)
			}
		}

		fmt.Println("Press 'Enter' to continue...")
		fmt.Scanln()
		rsp, err = chat.User("continue", nil, tools, "text", 0)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if strings.Contains(rsp.Choices[0].Message.Content, "FINISHED") {
			break
		}
		time.Sleep(time.Second)
	}
	rsp, err = chat.User("zdobyte DC_ID zwróć w postaci listy i nic poza listą nie zwracaj, przykład:[123,456]", nil, tools, "text", 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	fmt.Println()
	fmt.Println("Conversation:")
	fmt.Println(chat.DumpConversation())

	// Verify answer
	answerJSON := rsp.Choices[0].Message.Content
	var answer []int
	err = json.Unmarshal([]byte(answerJSON), &answer)
	if err != nil {
		log.Fatalf("Failed to deserialize answerJSON: %+v", err)
	}

	result, err := api.VerifyTaskAnswer("database", answer, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)

}
