package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// CancellationsService handles communication with the cancellation-related
// endpoints of the TrueTrial API.
type CancellationsService struct {
	client *httpClient
}

// Create initiates a cancellation for the given order.
func (s *CancellationsService) Create(ctx context.Context, orderID string, params CreateCancellationParams) (*Cancellation, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/cancellations", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data Cancellation `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode cancellation response: %w", err)
	}

	return &wrapper.Data, nil
}

// Get retrieves the cancellation for the given order.
func (s *CancellationsService) Get(ctx context.Context, orderID string) (*Cancellation, error) {
	body, err := s.client.get(ctx, "/orders/"+orderID+"/cancellations", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data Cancellation `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode cancellation response: %w", err)
	}

	return &wrapper.Data, nil
}
