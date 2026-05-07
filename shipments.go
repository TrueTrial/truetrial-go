package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// ShipmentsService handles communication with the shipment-related endpoints
// of the TrueTrial API.
type ShipmentsService struct {
	client *httpClient
}

// Create creates a new shipment for the given order.
func (s *ShipmentsService) Create(ctx context.Context, orderID string, params CreateShipmentParams) (*Shipment, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/shipments", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data Shipment `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode shipment response: %w", err)
	}

	return &wrapper.Data, nil
}

// List returns all shipments for the given order.
func (s *ShipmentsService) List(ctx context.Context, orderID string) ([]Shipment, error) {
	body, err := s.client.get(ctx, "/orders/"+orderID+"/shipments", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data []Shipment `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode shipments response: %w", err)
	}

	return wrapper.Data, nil
}
