package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nekitmilk/agent/internal/models"
)

type HTTPSender struct {
	baseURL string
	timeout time.Duration
	client  *http.Client
}

func NewHTTPSender(baseURL string, timeout time.Duration) *HTTPSender {
	return &HTTPSender{
		baseURL: baseURL,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *HTTPSender) SendMetrics(batch models.MetricsRequest) error {
	jsonData, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	url := fmt.Sprintf("%s/api/metrics", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
