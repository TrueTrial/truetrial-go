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

// ConfirmManuallyParams contains the fields for manually confirming delivery.
type ConfirmManuallyParams struct {
	DeliveredAt      string `json:"delivered_at"`
	Reason           string `json:"reason"`
	ConfirmedByEmail string `json:"confirmed_by_email,omitempty"`
}

// ConfirmManually manually confirms delivery of an order. Use this for the edge
// case where the carrier lost the package update but the consumer confirmed
// receipt. Records delivery_source = manual and starts the trial timer.
func (s *ShipmentsService) ConfirmManually(ctx context.Context, orderID string, params ConfirmManuallyParams) (*Order, error) {
	body, err := s.client.post(ctx, "/orders/"+orderID+"/confirm-delivery", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data Order `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode order response: %w", err)
	}

	return &wrapper.Data, nil
}
