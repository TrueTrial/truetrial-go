package truetrial

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeliveryFailedEnumValues(t *testing.T) {
	if WebhookEventDeliveryFailed != "delivery.failed" {
		t.Errorf("expected WebhookEventDeliveryFailed = 'delivery.failed', got %s", WebhookEventDeliveryFailed)
	}
	if OrderStatusDeliveryFailed != "delivery_failed" {
		t.Errorf("expected OrderStatusDeliveryFailed = 'delivery_failed', got %s", OrderStatusDeliveryFailed)
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.http.apiKey != "test-api-key" {
		t.Errorf("expected apiKey 'test-api-key', got '%s'", client.http.apiKey)
	}
	if client.http.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL '%s', got '%s'", DefaultBaseURL, client.http.baseURL)
	}
}

func TestNewClientWithBaseURL(t *testing.T) {
	client := NewClient("key", WithBaseURL("https://custom.example.com/api/v1"))
	if client.http.baseURL != "https://custom.example.com/api/v1" {
		t.Errorf("expected custom base URL, got '%s'", client.http.baseURL)
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 60 * time.Second}
	client := NewClient("key", WithHTTPClient(custom))
	if client.http.client != custom {
		t.Error("expected custom HTTP client to be set")
	}
}

func TestClientServicesInitialized(t *testing.T) {
	client := NewClient("key")
	if client.Orders == nil {
		t.Error("Orders service is nil")
	}
	if client.Shipments == nil {
		t.Error("Shipments service is nil")
	}
	if client.DigitalDelivery == nil {
		t.Error("DigitalDelivery service is nil")
	}
	if client.Temporal == nil {
		t.Error("Temporal service is nil")
	}
	if client.Cancellations == nil {
		t.Error("Cancellations service is nil")
	}
	if client.Webhooks == nil {
		t.Error("Webhooks service is nil")
	}
	if client.System == nil {
		t.Error("System service is nil")
	}
}

func TestAPIKeyHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Api-Key")
		if key != "test-key-123" {
			t.Errorf("expected X-Api-Key 'test-key-123', got '%s'", key)
		}
		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("expected Accept 'application/json', got '%s'", accept)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	}))
	defer server.Close()

	client := NewClient("test-key-123", WithBaseURL(server.URL))
	_, _ = client.Webhooks.List(context.Background())
}

func TestOrdersGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/orders/order-123" {
			t.Errorf("expected path '/orders/order-123', got '%s'", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":                "order-123",
				"tenant_id":         "tenant-1",
				"consumer_id":       "consumer-1",
				"external_order_id": "EXT-001",
				"product_name":      "Test Product",
				"product_type":      "physical",
				"product_price_cents": 2999,
				"product_currency":  "USD",
				"status":            "received",
				"created_at":        "2026-01-15T10:00:00Z",
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	order, err := client.Orders.Get(context.Background(), "order-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != "order-123" {
		t.Errorf("expected order ID 'order-123', got '%s'", order.ID)
	}
	if order.ProductPriceCents != 2999 {
		t.Errorf("expected price 2999, got %d", order.ProductPriceCents)
	}
	if order.Status != OrderStatusReceived {
		t.Errorf("expected status 'received', got '%s'", order.Status)
	}
}

func TestOrdersList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got '%s'", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("status") != "delivered" {
			t.Errorf("expected status=delivered, got '%s'", r.URL.Query().Get("status"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"data":         []any{},
			"current_page": 2,
			"last_page":    5,
			"per_page":     15,
			"total":        72,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.Orders.List(context.Background(), &ListOrdersParams{
		Page:   2,
		Status: OrderStatusDelivered,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CurrentPage != 2 {
		t.Errorf("expected current page 2, got %d", resp.CurrentPage)
	}
	if resp.Total != 72 {
		t.Errorf("expected total 72, got %d", resp.Total)
	}
}

func TestErrorHandling401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API key"})
	}))
	defer server.Close()

	client := NewClient("bad-key", WithBaseURL(server.URL))
	_, err := client.Orders.Get(context.Background(), "order-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsAuthenticationError(err) {
		t.Errorf("expected authentication error, got: %v", err)
	}
}

func TestErrorHandling404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order not found"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.Orders.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFoundError(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestErrorHandling422(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "The given data was invalid.",
			"errors": map[string][]string{
				"product_name": {"The product name field is required."},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.Orders.Create(context.Background(), CreateOrderParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsValidationError(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	ttErr := err.(*TrueTrialError)
	if len(ttErr.Errors["product_name"]) == 0 {
		t.Error("expected product_name validation error")
	}
}

func TestErrorHandling429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"message": "Too many requests"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.Orders.Get(context.Background(), "order-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimitError(err) {
		t.Errorf("expected rate limit error, got: %v", err)
	}
	ttErr := err.(*TrueTrialError)
	if ttErr.RetryAfter != 30 {
		t.Errorf("expected RetryAfter 30, got %d", ttErr.RetryAfter)
	}
}

func TestErrorHandling500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.Orders.Get(context.Background(), "order-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsServerError(err) {
		t.Errorf("expected server error, got: %v", err)
	}
}

func TestVerifyWebhookSignature(t *testing.T) {
	payload := []byte(`{"event":"order.created","data":{"id":"order-123"}}`)
	secret := "whsec_test_secret_key"
	signature := computeHMAC(payload, secret)

	if !VerifyWebhookSignature(payload, signature, secret) {
		t.Error("expected valid signature to pass verification")
	}
}

func TestVerifyWebhookSignatureInvalid(t *testing.T) {
	payload := []byte(`{"event":"order.created","data":{"id":"order-123"}}`)
	secret := "whsec_test_secret_key"

	if VerifyWebhookSignature(payload, "invalid-signature", secret) {
		t.Error("expected invalid signature to fail verification")
	}
}

func TestVerifyWebhookSignatureTampered(t *testing.T) {
	payload := []byte(`{"event":"order.created","data":{"id":"order-123"}}`)
	tampered := []byte(`{"event":"order.created","data":{"id":"order-999"}}`)
	secret := "whsec_test_secret_key"
	signature := computeHMAC(payload, secret)

	if VerifyWebhookSignature(tampered, signature, secret) {
		t.Error("expected tampered payload to fail verification")
	}
}

func TestVerifyWebhookSignatureWithTolerance(t *testing.T) {
	payload := []byte(`{"event":"order.created","data":{"id":"order-123"}}`)
	secret := "whsec_test_secret_key"
	ts := fmt.Sprintf("%d", time.Now().Unix())

	signedPayload := fmt.Sprintf("%s.%s", ts, string(payload))
	hash := computeHMAC([]byte(signedPayload), secret)
	signature := fmt.Sprintf("t=%s,v1=%s", ts, hash)

	if !VerifyWebhookSignatureWithTolerance(payload, signature, secret, 5*time.Minute) {
		t.Error("expected valid signature with tolerance to pass")
	}
}

func TestVerifyWebhookSignatureWithToleranceExpired(t *testing.T) {
	payload := []byte(`{"event":"order.created","data":{"id":"order-123"}}`)
	secret := "whsec_test_secret_key"
	// Timestamp 10 minutes in the past.
	ts := fmt.Sprintf("%d", time.Now().Add(-10*time.Minute).Unix())

	signedPayload := fmt.Sprintf("%s.%s", ts, string(payload))
	hash := computeHMAC([]byte(signedPayload), secret)
	signature := fmt.Sprintf("t=%s,v1=%s", ts, hash)

	if VerifyWebhookSignatureWithTolerance(payload, signature, secret, 5*time.Minute) {
		t.Error("expected expired signature to fail tolerance check")
	}
}

func TestVerifyWebhookSignatureWithToleranceMalformed(t *testing.T) {
	payload := []byte(`{"event":"order.created"}`)
	secret := "whsec_test_secret_key"

	if VerifyWebhookSignatureWithTolerance(payload, "garbage", secret, 5*time.Minute) {
		t.Error("expected malformed signature to fail")
	}
	if VerifyWebhookSignatureWithTolerance(payload, "t=abc,v1=def", secret, 5*time.Minute) {
		t.Error("expected non-numeric timestamp to fail")
	}
}

func TestWebhooksDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/webhooks/wh-123" {
			t.Errorf("expected path '/webhooks/wh-123', got '%s'", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.Webhooks.Delete(context.Background(), "wh-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShipmentsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/orders/order-1/shipments" {
			t.Errorf("expected path '/orders/order-1/shipments', got '%s'", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":              "ship-1",
				"order_id":        "order-1",
				"carrier":         "ups",
				"tracking_number": "1Z999AA10123456784",
				"status":          "pending",
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	shipment, err := client.Shipments.Create(context.Background(), "order-1", CreateShipmentParams{
		Carrier:        CarrierUPS,
		TrackingNumber: "1Z999AA10123456784",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipment.ID != "ship-1" {
		t.Errorf("expected shipment ID 'ship-1', got '%s'", shipment.ID)
	}
	if shipment.Carrier != CarrierUPS {
		t.Errorf("expected carrier 'ups', got '%s'", shipment.Carrier)
	}
}

func TestEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"OrderStatusReceived", string(OrderStatusReceived), "received"},
		{"OrderStatusCancelled", string(OrderStatusCancelled), "cancelled"},
		{"TemporalTypeTrial", string(TemporalTypeTrial), "trial"},
		{"TemporalTypeGuarantee", string(TemporalTypeGuarantee), "guarantee"},
		{"TemporalStatusClaimApproved", string(TemporalStatusClaimApproved), "claim_approved"},
		{"ShipmentStatusOutForDelivery", string(ShipmentStatusOutForDelivery), "out_for_delivery"},
		{"ProductTypeDigital", string(ProductTypeDigital), "digital"},
		{"CarrierFedEx", string(CarrierFedEx), "fedex"},
		{"DurationUnitMonths", string(DurationUnitMonths), "months"},
		{"DeliverySourceWebhook", string(DeliverySourceWebhook), "webhook"},
		{"WebhookEventDisputeLost", string(WebhookEventDisputeLost), "dispute.lost"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, tt.value)
			}
		})
	}
}
