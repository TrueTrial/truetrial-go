package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// WebhooksService handles communication with the webhook subscription
// endpoints of the TrueTrial API.
type WebhooksService struct {
	client *httpClient
}

// List returns all webhook subscriptions for the authenticated tenant.
func (s *WebhooksService) List(ctx context.Context) ([]WebhookSubscription, error) {
	body, err := s.client.get(ctx, "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data []WebhookSubscription `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode webhooks response: %w", err)
	}

	return wrapper.Data, nil
}

// Create registers a new webhook subscription and returns it.
// The response includes the signing secret which should be stored securely.
func (s *WebhooksService) Create(ctx context.Context, params CreateWebhookParams) (*WebhookSubscription, error) {
	body, err := s.client.post(ctx, "/webhooks", params)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data WebhookSubscription `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode webhook response: %w", err)
	}

	return &wrapper.Data, nil
}

// Delete removes a webhook subscription by its ID.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	_, err := s.client.delete(ctx, "/webhooks/"+id)
	return err
}
