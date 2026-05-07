package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// OrdersService handles communication with the order-related endpoints
// of the TrueTrial API.
type OrdersService struct {
	client *httpClient
}

// List returns a paginated list of orders. Pass nil for params to use defaults.
func (s *OrdersService) List(ctx context.Context, params *ListOrdersParams) (*PaginatedResponse, error) {
	query := url.Values{}
	if params != nil {
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.Status != "" {
			query.Set("status", string(params.Status))
		}
	}

	body, err := s.client.get(ctx, "/orders", query)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode orders response: %w", err)
	}

	return &resp, nil
}

// Create creates a new order and returns it.
func (s *OrdersService) Create(ctx context.Context, params CreateOrderParams) (*Order, error) {
	body, err := s.client.post(ctx, "/orders", params)
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

// Get retrieves a single order by its ID.
func (s *OrdersService) Get(ctx context.Context, id string) (*Order, error) {
	body, err := s.client.get(ctx, "/orders/"+id, nil)
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

// Status retrieves the combined status of an order, its temporal element,
// and its shipment.
func (s *OrdersService) Status(ctx context.Context, id string) (*OrderStatusResponse, error) {
	body, err := s.client.get(ctx, "/orders/"+id+"/status", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data OrderStatusResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode order status response: %w", err)
	}

	return &wrapper.Data, nil
}
