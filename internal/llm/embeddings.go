package llm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	backoffPeriod   = 200 * time.Millisecond
	embedMaxRetries = 3
)

type embeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Embed generates an embedding for the given text.
//
// It retries on transient failures (network errors and 5xx responses)
// using a simple backoff strategy. Non-retryable errors (e.g. 4xx)
// are returned immediately.
func (c *Client) Embed(text string) ([]float64, error) {
	reqBody := embeddingRequest{
		Prompt: text,
		Model:  c.EmbeddingsModel,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	for i := range embedMaxRetries {
		resp, err := c.post(embeddingsEndpoint, data)
		if err != nil {
			// network error -> retry
			time.Sleep(time.Duration(i+1) * backoffPeriod)
			continue
		}

		if resp.StatusCode >= 500 {
			// ollama unavailable
			resp.Body.Close()
			time.Sleep(time.Duration(i+1) * backoffPeriod)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("llm returned status %d", resp.StatusCode)
		}

		var res embeddingResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		return res.Embedding, nil
	}

	return nil, fmt.Errorf("embed failed after retries")
}
