package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ErrRateLimit = fmt.Errorf("rate limit exceeded")
)

type AI interface {
	Query(query string) (string, error)
}

type ai struct {
	client http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QueryRequest struct {
	Model    string    `json:"model"`
	Store    bool      `json:"store"`
	Messages []Message `json:"messages"`
}

func New() AI {
	cli := http.Client{}

	return &ai{
		client: cli,
	}
}

func (a *ai) Query(query string) (string, error) {

	url := "https://api.openai.com/v1/chat/completions"

	q := QueryRequest{
		Model: "gpt-3.5-turbo",
		Store: true,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	data, err := json.Marshal(q)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}

	dataResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	dataRespStr := string(dataResp)
	dataRespStr = strings.TrimPrefix(dataRespStr, "```json ")
	dataRespStr = strings.TrimSuffix(dataRespStr, " ```")

	if resp.StatusCode == 429 {
		fmt.Println(dataRespStr)
		return "", ErrRateLimit
	}

	if os.Getenv("DEBUG") == "1" {
		fmt.Println(dataRespStr)
	}

	return string(dataResp), nil
}
