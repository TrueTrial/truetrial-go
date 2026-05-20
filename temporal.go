package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// TemporalService handles communication with the temporal element endpoints
// of the TrueTrial API.
type TemporalService struct {
	client *httpClient
}

// Get retrieves the temporal element for the given order.
func (s *TemporalService) Get(ctx context.Context, orderID string) (*TemporalElement, error) {
	body, err := s.client.get(ctx, "/orders/"+orderID+"/temporal", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data TemporalElement `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode temporal response: %w", err)
	}

	return &wrapper.Data, nil
}

// Extend extends the temporal element for the given order by an additional
// duration.
func (s *TemporalService) Extend(ctx context.Context, orderID string, params ExtendTemporalParams) (*TemporalElement, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/temporal/extend", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data TemporalElement `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode temporal response: %w", err)
	}

	return &wrapper.Data, nil
}

// Adjust sets a new end time for the temporal element on the given order.
func (s *TemporalService) Adjust(ctx context.Context, orderID string, params AdjustTemporalParams) (*TemporalElement, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/temporal/adjust", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data TemporalElement `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode temporal response: %w", err)
	}

	return &wrapper.Data, nil
}

// Claim submits a warranty or guarantee claim for the given order.
func (s *TemporalService) Claim(ctx context.Context, orderID string, params ClaimParams) (*TemporalElement, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/temporal/claim", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data TemporalElement `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode temporal response: %w", err)
	}

	return &wrapper.Data, nil
}

// ResolveClaim resolves an existing claim on the given order.
func (s *TemporalService) ResolveClaim(ctx context.Context, orderID string, params ResolveClaimParams) (*TemporalElement, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/temporal/resolve-claim", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data TemporalElement `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode temporal response: %w", err)
	}

	return &wrapper.Data, nil
}
