package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/api"
	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

/*
Scenario:
- get image urls as a list (tool: ask api)
- ask about each condition of each
- fix each image based on it's condition
- describe the person on every image (if any person visible)
- return collective detailed description of the person

Tools:
- get_images: unstructured string
- assess_condition(url): TOO-DARK/TOO-BRIGHT/DAMAGE
- make_brighter(url): url
- make_darker(url): url
- fix_damaged(url): url
- describe_image(url)
*/
func main() {
	const system = "Rozwiąz zadanie krok po kroku, korzystając z dostępnych narzędzi gdzie to tylko mozliwe. Nie zaczynaj od planu działania; po prostu zacznij wykonywać kroki."
	const user = `Masz stworzyć krótki i zwięzły rysopis kobiety widocznej na zadanych obrazach(zdjęciach). Kroki:
	1. Pobierz informację o lokalizacji obrazów i przekształć ją do postaci listy URL. Przykład listy URL: [https://wiki.com/apple.jpg,https://wiki.com/pear.jpg].
	2. Przekształć URL kazdego z obrazów dodając suffix "-small". Przykład: https://wiki.com/apple.jpg -> https://wiki.com/apple-small.jpg
	3. Sprawdź stan każdego obrazu z uzyciem właściwego narzędzia.
	4. Popraw każdy obraz w oparciu o jego stan z uzyciem właściwego narzędzia.
	5. Wszystkie poprawione obrazy przedstawiają tę samą kobietę. Po kolei zapytaj o opis osoby przedstawionej na każdym z obrazów.
	7. Na koniec napisz: [FINISHED]
	Do dzieła, wykonaj kroki!
	`
	// user := "sprawdź stan kazdego z obrazów, korzystając z odpowiedniego narzędzia: https://wiki.com/apple.jpg, https://wiki.com/pear.jpg"
	const debug = true

	tools := prepareTools()
	chat, err := openai.NewChatWithMemory(system, "gpt-4o", 1000, debug)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	rsp, err := chat.User(user, nil, tools, "text", 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// repeat until [FINISHED]
FINISHED:
	for {
		for _, choice := range rsp.Choices {
			if choice.Message.Content != "" {
				if strings.Contains(choice.Message.Content, "FINISHED") {
					break FINISHED
				}
			}

			responses := []openai.ToolCallResponse{}
			for _, toolCall := range choice.Message.ToolCalls {
				name := toolCall.Function.Name
				args := getArgs(toolCall.Function.Arguments)
				fmt.Println("TOOL:", name, args)
				callID := toolCall.ID
				toolResponse := ""
				switch name {
				case "get_image_urls":
					toolResponse = askAPI("START")
				case "check_image_condition":
					toolResponse = checkImageCondition(args["image_url"])
				case "repair_image":
					toolResponse = askAPI("REPAIR " + filepath.Base(args["image_url"]))
				case "brighten_image":
					toolResponse = askAPI("BRIGHTEN " + filepath.Base(args["image_url"]))
				case "describe_image":
					toolResponse = describeImage(args["image_url"])
				default:
					log.Fatal("unknown function: " + name)
				}
				responses = append(responses, openai.ToolCallResponse{Text: toolResponse, ToolCallID: callID})
			}
			rsp, err = chat.ToolResponse(responses)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			if strings.Contains(rsp.Choices[0].Message.Content, "FINISHED") {
				break FINISHED
			}
		}

		fmt.Println("Press 'Enter' to continue...")
		fmt.Scanln()
		rsp, err = chat.User("continue", nil, tools, "text", 0)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		time.Sleep(time.Second) // protect from infinite loop of llm requests :)
	}

	rsp, err = chat.User("Połącz opisy obrazów w jeden krótki, spójny opis kobiety. Nie pisz juz [FINISHED]", nil, tools, "text", 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	result, err := api.VerifyTaskAnswer("photos", rsp, api.VerificationURL)
	if err != nil {
		fmt.Println("Answer verification failed:", err)
		return
	}
	fmt.Println(result)
}

func checkImageCondition(url string) string {
	system := "you are image condition assesor expert"
	user := "check the condition of attached image and return the condition as either: TOO-BRIGHT or TOO-DARK or DAMAGED or GOOD"
	images := []string{openai.ImageFromURL(url)}
	response, err := openai.CompletionCheap(user, system, images)
	if err != nil {
		return err.Error()
	}
	return response
}
func describeImage(url string) string {
	system := "you are image written descriptions expert"
	user := "describe the people in the image, if there are people. If there are no people - return empty response"
	images := []string{openai.ImageFromURL(url)}
	response, err := openai.CompletionCheap(user, system, images)
	if err != nil {
		return err.Error()
	}
	return response
}

func getArgs(argsJSON string) map[string]string {
	var args map[string]string
	err := json.Unmarshal([]byte(argsJSON), &args)
	if err != nil {
		log.Fatalf("failed to unmarshal argsJSON: %+v", err)
	}
	return args
}

// prepare functions for GPT to select from
func prepareTools() []openai.Tool {
	// prepare temp function parameters
	getImageConditionFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"image_url": {
				Type:        "string",
				Description: "URL of the image to assess the condition, e.g. 'https://example.com/ball.jpg'",
			},
		},
	}
	repairImageFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"image_url": {
				Type:        "string",
				Description: "URL of the image to be repaired, e.g. 'https://example.com/ball.jpg'",
			},
		},
	}
	brightenImageFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"image_url": {
				Type:        "string",
				Description: "URL of the image to be brightened, e.g. 'https://example.com/ball.jpg'",
			},
		},
	}
	describeImageFuncParams := openai.Parameters{
		Type: "object",
		Properties: map[string]openai.Property{
			"image_url": {
				Type:        "string",
				Description: "URL of the image to be described, e.g. 'https://example.com/ball.jpg'",
			},
		},
	}

	// prepare functions
	checkImageConditionsFunc := openai.Function{
		Name:        "check_image_condition",
		Description: "this function checks the condition of given image and returns the condition as: TOO-BRIGHT or TOO-DARK or DAMAGED or GOOD",
		Parameters:  &getImageConditionFuncParams,
	}
	getImageURLsFunc := openai.Function{
		Name:        "get_image_urls",
		Description: "this function returns descriptive information about the location of files",
		Parameters:  nil, // no parameters
	}
	repairImageFunc := openai.Function{
		Name:        "repair_image",
		Description: "this function repairs the image of given URL and returns the resulting image URL",
		Parameters:  &repairImageFuncParams,
	}
	brightenImageFunc := openai.Function{
		Name:        "brighten_image",
		Description: "this function brightens the image of given URL and returns the resulting image URL",
		Parameters:  &brightenImageFuncParams,
	}
	describeImageFunc := openai.Function{
		Name:        "describe_image",
		Description: "this function describes the image of given URL and returns the description",
		Parameters:  &describeImageFuncParams,
	}

	// prepare tools from functions
	getImageUrls := openai.Tool{
		Type:     "function",
		Function: getImageURLsFunc,
	}
	checkImage := openai.Tool{
		Type:     "function",
		Function: checkImageConditionsFunc,
	}
	repairImage := openai.Tool{
		Type:     "function",
		Function: repairImageFunc,
	}
	brightenImage := openai.Tool{
		Type:     "function",
		Function: brightenImageFunc,
	}
	describeImage := openai.Tool{
		Type:     "function",
		Function: describeImageFunc,
	}

	// return tool set
	return []openai.Tool{getImageUrls, repairImage, brightenImage, checkImage, describeImage}
}

func askAPI(question string) string {
	response, err := api.VerifyTaskAnswer("photos", question, api.VerificationURL)
	if err != nil {
		return err.Error()
	}
	return response
}
