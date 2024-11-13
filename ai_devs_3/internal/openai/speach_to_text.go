package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const speachToTextURL = "https://api.openai.com/v1/audio/transcriptions"

type GPTErrorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    *string `json:"code"`
	} `json:"error"`
}

// input audio file size limit is 25MB
func SpeachToText(inputMP3 string) (string, error) {
	// get API KEY from env
	apiKey := os.Getenv("OPENAI_API_KEY") // Get the API key from the environment variable
	if apiKey == "" {
		return "", errors.New("OpenAI API key is not set")
	}

	// open input audio file
	file, err := os.Open(inputMP3)
	if err != nil {
		return "", errors.Wrap(err, "failed to open input audio file")
	}
	defer file.Close()

	// send request
	req, err := prepareSpeachToTextRequest(file, apiKey)
	if err != nil {
		return "", errors.WithMessage(err, "failed to prepare speech to text request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != http.StatusOK {
		var errorResponse GPTErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return "", errors.Wrap(err, "failed to unmarshal error response")
		}
		errorMsg := fmt.Sprintf("request failed [code:%d]: %s [type:%s]", resp.StatusCode, errorResponse.Error.Message, errorResponse.Error.Type)
		return "", errors.New(errorMsg)
	}

	// return transcription
	return string(body), nil
}

// https://platform.openai.com/docs/api-reference/audio/createTranscription
func prepareSpeachToTextRequest(file *os.File, apiKey string) (*http.Request, error) {
	// prepare request body
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	// attach the audio file itself
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create form file for audio")
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to copy audio file to form")
	}

	// attach model
	model, err := writer.CreateFormField("model")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create form field for model")
	}
	_, err = model.Write([]byte("whisper-1")) // [whisper-1]
	if err != nil {
		return nil, errors.Wrap(err, "failed to write model to form field")
	}

	// attach response_format
	format, err := writer.CreateFormField("response_format")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create form field for response format")
	}
	_, err = format.Write([]byte("text")) // [json, text, srt, verbose_json, vtt]
	if err != nil {
		return nil, errors.Wrap(err, "failed to write response format to form field")
	}

	// flush to buffer
	writer.Close()

	// compose the request
	req, err := http.NewRequest("POST", speachToTextURL, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new HTTP request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// SUCCESS
	return req, nil
}
