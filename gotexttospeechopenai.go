// 15$ = 1mln characters

package texttospeechopenai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

const URL = "https://api.openai.com/v1/audio/speech"

var allowedVoices = []string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}

var allowedModels = []string{"tts1", "tts1hd"}

type CustomError struct {
	Message string
}

type Input struct {
	model string `json:"model"`
	voice string `json:"voice"`
	text  string `json:"text"`
}

type TextToSpeechClient struct {
	apiParams
}

type apiParams struct {
	Client *http.Client
	APIKey string
}

func NewTextToSpeechClient(apiParams apiParams) *TextToSpeechClient {
	return &TextToSpeechClient{
		apiParams: apiParams,
	}
}

func (e *CustomError) Error() string {
	return e.Message
}

func isValidateInput(input string, allowed []string) bool {
	for _, val := range allowed {
		if val == input {
			return true
		}
	}

	return false
}

func (c *TextToSpeechClient) TextToSpeech(input Input, filePath string) error {

	if input.text == "" {
		return &CustomError{Message: fmt.Sprintf("empty text provided: %s", input.text)}
	}

	if !isValidateInput(input.model, allowedModels) {
		return &CustomError{Message: fmt.Sprintf("invalid input model: %s. allowed values are: %v", input, allowedModels)}
	}

	if !isValidateInput(input.voice, allowedVoices) {
		return &CustomError{Message: fmt.Sprintf("invalid input voice: %s. allowed values are: %v", input.voice, allowedVoices)}
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return &CustomError{Message: fmt.Sprintf("error marshalling to JSON: %v", err)}
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiParams.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.apiParams.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	outFile, err := os.Create(path.Join(filePath, input.text[:30]))
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, req.Body)
	if err != nil {
		return err
	}
	return nil
}
