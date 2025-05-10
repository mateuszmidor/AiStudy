package api

import (
	"encoding/json"
	"log"
)

func FillDataFromJSONURL(url string, dst any) {
	jsonStr, err := FetchData(url)
	if err != nil {
		log.Fatal(err)
	}
	FillDataFromJSON(jsonStr, dst)
}

func FillDataFromJSON(jsonStr string, dst any) {
	err := json.Unmarshal([]byte(jsonStr), dst)
	if err != nil {
		log.Fatal(err)
	}
}
