package ollama

import (
	"bytes"
	"context"
	"core_service/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL    string
	Model      string
	httpClient *http.Client
}

func InitOllama(cfg config.OllamaConfig) *Client {
	return &Client{
		BaseURL: cfg.BaseUrl,
		Model:   cfg.Model,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/api/generate",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var response GenerateResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Response, nil
}
