package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// DigitalDeliveryService handles communication with the digital delivery
// endpoints of the TrueTrial API.
type DigitalDeliveryService struct {
	client *httpClient
}

// Confirm marks a digital product as delivered for the given order.
func (s *DigitalDeliveryService) Confirm(ctx context.Context, orderID string, params ConfirmDigitalDeliveryParams) (*Order, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/digital-delivery", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data Order `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode digital delivery response: %w", err)
	}

	return &wrapper.Data, nil
}
