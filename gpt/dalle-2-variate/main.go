package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const openAIURL = "https://api.openai.com/v1/images/variations"

type DalleGenerateResponse struct {
	Created int64 `json:"created"` // unix timestamp
	Data    []struct {
		URL         string `json:"url"`      // either this
		Base64Image string `json:"b64_json"` // or this
	} `json:"data"`
}

type DalleErrorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// must be a valid PNG file, less than 4MB, and square
func variateImage(inputPNG string) string {
	// get API KEY from env
	apiKey := os.Getenv("GPT_APIKEY") // Get the API key from the environment variable
	if apiKey == "" {
		panic("OpenAI API key is not set")
	}

	// open input PNG file
	file, err := os.Open(inputPNG)
	panicOnError(err)
	defer file.Close()

	// send request
	req := prepareVariateImageRequest(file, apiKey)
	resp, err := http.DefaultClient.Do(req)
	panicOnError(err)
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	panicOnError(err)
	if resp.StatusCode != http.StatusOK {
		var errorResponse DalleErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		panicOnError(err)
		errorMsg := fmt.Sprintf("request failed [code:%d]: %s [type:%s]", resp.StatusCode, errorResponse.Error.Message, errorResponse.Error.Type)
		panic(errorMsg)
	}

	// deserialize response
	var generateResponse DalleGenerateResponse
	err = json.Unmarshal(body, &generateResponse)
	panicOnError(err)

	// assuming you want to return the URL of the first generated image
	if len(generateResponse.Data) > 0 {
		return generateResponse.Data[0].URL
	}

	return "no image URL returned"
}

// https://platform.openai.com/docs/api-reference/images/createVariation
func prepareVariateImageRequest(file *os.File, apiKey string) *http.Request {
	// prepare request body
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	// attach the image file itself
	part, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
	panicOnError(err)
	_, err = io.Copy(part, file)
	panicOnError(err)

	// attach size
	size, err := writer.CreateFormField("size")
	panicOnError(err)
	_, err = size.Write([]byte("1024x1024")) // [256x256, 512x512, 1024x1024]
	panicOnError(err)

	// attach num of variations to generate
	numVariations, err := writer.CreateFormField("n")
	panicOnError(err)
	_, err = numVariations.Write([]byte("1"))
	panicOnError(err)

	// attach response format
	responseFormat, err := writer.CreateFormField("response_format")
	panicOnError(err)
	_, err = responseFormat.Write([]byte("url")) // [url, b64_json]
	panicOnError(err)

	// flush to buffer
	writer.Close()

	// compose the request
	req, err := http.NewRequest("POST", openAIURL, reqBody)
	panicOnError(err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req
}

// user panicOnError to reduce lines of code by 50%.
func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	const imagePNG = "corgi.png"

	url := variateImage(imagePNG)
	fmt.Println(url)
}
