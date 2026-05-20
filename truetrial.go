// Package truetrial provides a Go client for the TrueTrial API.
//
// TrueTrial is a compliance-first platform for managing trial, warranty,
// subscription, and guarantee periods. This SDK provides full access to
// the TrueTrial REST API and utilities for webhook signature verification.
//
// Usage:
//
//	client := truetrial.NewClient("your-api-key")
//	order, err := client.Orders.Get(ctx, "order-id")
package truetrial

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default API base URL.
	DefaultBaseURL = "https://truetrial.test/api/v1"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// Version is the SDK version.
	Version = "1.0.0"
)

// Client is the TrueTrial API client. Use NewClient to create one.
type Client struct {
	Orders          *OrdersService
	Shipments       *ShipmentsService
	DigitalDelivery *DigitalDeliveryService
	Temporal        *TemporalService
	Cancellations   *CancellationsService
	Webhooks        *WebhooksService
	System          *SystemService

	http *httpClient
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.http.baseURL = url
	}
}

// WithHTTPClient overrides the default net/http client used for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.http.client = hc
	}
}

// NewClient creates a new TrueTrial API client authenticated with the given
// API key. Options can be provided to customize the client behaviour.
func NewClient(apiKey string, opts ...Option) *Client {
	hc := &httpClient{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	c := &Client{http: hc}

	for _, opt := range opts {
		opt(c)
	}

	c.Orders = &OrdersService{client: hc}
	c.Shipments = &ShipmentsService{client: hc}
	c.DigitalDelivery = &DigitalDeliveryService{client: hc}
	c.Temporal = &TemporalService{client: hc}
	c.Cancellations = &CancellationsService{client: hc}
	c.Webhooks = &WebhooksService{client: hc}
	c.System = &SystemService{client: hc}

	return c
}
