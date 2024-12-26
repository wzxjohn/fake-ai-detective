package main

import (
	"net/http"
	"time"
)

type TrackedResult struct {
	Timestamp int64       `json:"timestamp,omitempty"`
	IP        string      `json:"ip,omitempty"`
	UserAgent string      `json:"user_agent,omitempty"`
	Message   string      `json:"message,omitempty"`
	Header    http.Header `json:"header,omitempty"`
}

// TrackedRequest 存储请求的相关信息
type TrackedRequest struct {
	Timestamp time.Time
	Image     []byte
	Finished  bool
	Results   []*TrackedResult
}

// OpenAIRequest 结构体用于构造发送到 OpenAI 的请求
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type StartRequest struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

type StartResponse APIResponse[string]

type DetectResult struct {
	Finished bool             `json:"finished"`
	Results  []*TrackedResult `json:"results"`
}
type ResultResponse APIResponse[*DetectResult]
