package health

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HealthService interface {
	Check(ctx context.Context) error
}

type HTTPHealthService struct {
	URL     string
	Timeout time.Duration
}

func NewHTTPHealthService(url string, timeout time.Duration) *HTTPHealthService {
	return &HTTPHealthService{
		URL:     url,
		Timeout: timeout,
	}
}

func (h *HTTPHealthService) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("health request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("could not close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
