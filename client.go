package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	client = http.Client{Timeout: 120 * time.Second}
)

func sendOpenAIRequest(url, key, model string, traceID string) {
	defer markFinished(traceID)
	headers := map[string]string{
		"Accept":        "",
		"User-Agent":    "Apifox/1.0.0 (https://apifox.com)",
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + key,
	}

	imageURL := baseURL + traceID
	requestBody := OpenAIRequest{
		Model: model,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "image_url",
						ImageURL: &ImageURL{
							URL: imageURL,
						},
					},
					{
						Type: "text",
						Text: "What is this?",
					},
				},
			},
		},
		MaxTokens: 3,
		Stream:    false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		recordMessage(traceID, time.Now().Unix(), "", "", fmt.Sprintf("Error: %v", err), nil)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		recordMessage(traceID, time.Now().Unix(), "", "", fmt.Sprintf("Error: %v", err), nil)
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		recordMessage(traceID, time.Now().Unix(), "", "", fmt.Sprintf("Exception: %v", err), nil)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		recordMessage(traceID, time.Now().Unix(), "", "", fmt.Sprintf("Error: Status %d", resp.StatusCode), nil)
	} else {
		recordMessage(traceID, time.Now().Unix(), "", "", "API 响应结束，完成探测", nil)
	}
}
