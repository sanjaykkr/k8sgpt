package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const customRAGAIClientName = "nova"

type CustomRagAIClient struct {
	nopCloser
}

func (c *CustomRagAIClient) Configure(_ IAIConfig) error {
	return nil
}

type RagResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func (c *CustomRagAIClient) GetCompletion(_ context.Context, prompt string) (string, error) {
	event := map[string]interface{}{
		"event": map[string]string{
			"channel":   "string",
			"email":     "string",
			"text":      prompt,
			"thread_ts": "string",
			"ts":        "string",
			"user":      "string",
		},
	}

	// Marshal event into JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		return "error happened", errors.New("error marshaling JSON data")
	}

	// Print the JSON data being sent
	fmt.Println("JSON Data Sent:", string(jsonData))

	resp, err := http.Post("https://nova-api.staging-service.nr-ops.net/query", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "error happened", errors.New("error getting response from custom RAG client")
	}
	defer resp.Body.Close()

	// Read and print the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error happened", errors.New("error reading response body")
	}

	fmt.Println("Response Body:", string(bodyBytes))

	var res RagResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "error happened", errors.New("error unmarshaling response")
	}

	return res.Body, nil
}

func (c *CustomRagAIClient) GetName() string {
	return customRAGAIClientName
}
