package main

import (
	"context"
	"fmt"
	"log"

	baml "baml/baml_client"
)

func main() {
	cfg := baml.NewConfiguration()
	b := baml.NewAPIClient(cfg).DefaultAPI // The BAML preview server runs on http://localhost:2024 (default)

	req := baml.NewGetRandomCityRequest("China")
	resp, r, err := b.GetRandomCity(context.Background()).GetRandomCityRequest(*req).Execute()
	if err != nil {
		fmt.Printf("Error when calling GetRandomCity: %v\n", err)
		fmt.Printf("Full HTTP response: %v\n", r)
		return
	}
	log.Printf("Random city: %s in %s", resp.City, resp.Country)
}
