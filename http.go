package truetrial

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// httpClient handles low-level HTTP communication with the TrueTrial API.
type httpClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// get performs an authenticated GET request. query may be nil.
func (h *httpClient) get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	fullURL := h.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("truetrial: failed to create request: %w", err)
	}

	return h.do(req)
}

// post performs an authenticated POST request with a JSON body.
func (h *httpClient) post(ctx context.Context, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("truetrial: failed to encode request body: %w", err)
		}
		reader = bytes.NewReader(encoded)
	}

	fullURL := h.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, reader)
	if err != nil {
		return nil, fmt.Errorf("truetrial: failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return h.do(req)
}

// delete performs an authenticated DELETE request.
func (h *httpClient) delete(ctx context.Context, path string) ([]byte, error) {
	fullURL := h.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("truetrial: failed to create request: %w", err)
	}

	return h.do(req)
}

// do executes the request with common headers and error handling.
func (h *httpClient) do(req *http.Request) ([]byte, error) {
	req.Header.Set("X-Api-Key", h.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "truetrial-go/"+Version)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("truetrial: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("truetrial: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return respBody, nil
	}

	return nil, parseErrorResponse(resp.StatusCode, resp.Header, respBody)
}

// parseErrorResponse creates a typed error from an HTTP error response.
func parseErrorResponse(statusCode int, headers http.Header, body []byte) error {
	// Attempt to extract a message from the JSON response body.
	var parsed struct {
		Message string            `json:"message"`
		Errors  map[string][]string `json:"errors"`
	}
	_ = json.Unmarshal(body, &parsed)

	msg := parsed.Message
	if msg == "" {
		msg = http.StatusText(statusCode)
	}

	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return NewAuthenticationError(msg, body)

	case http.StatusUnprocessableEntity:
		return NewValidationError(msg, parsed.Errors, body)

	case http.StatusNotFound:
		return NewNotFoundError(msg, body)

	case http.StatusTooManyRequests:
		retryAfter := 0
		if ra := headers.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return NewRateLimitError(msg, retryAfter, body)

	default:
		return NewServerError(statusCode, msg, body)
	}
}
