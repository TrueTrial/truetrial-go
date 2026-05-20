package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
)

// SystemService handles communication with system-level endpoints
// of the TrueTrial API.
type SystemService struct {
	client *httpClient
}

// CarrierHealth returns the health status of all configured shipping carriers.
func (s *SystemService) CarrierHealth(ctx context.Context) ([]CarrierHealthEntry, error) {
	body, err := s.client.get(ctx, "/carrier-health", nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data []CarrierHealthEntry `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("truetrial: failed to decode carrier health response: %w", err)
	}

	return wrapper.Data, nil
}
