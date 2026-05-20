package truetrial

import "encoding/json"

// ---------------------------------------------------------------------------
// Enums
// ---------------------------------------------------------------------------

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusReceived       OrderStatus = "received"
	OrderStatusShipped        OrderStatus = "shipped"
	OrderStatusInTransit      OrderStatus = "in_transit"
	OrderStatusDelivered      OrderStatus = "delivered"
	OrderStatusDeliveryFailed OrderStatus = "delivery_failed"
	OrderStatusTrialActive    OrderStatus = "trial_active"
	OrderStatusConverted      OrderStatus = "converted"
	OrderStatusReturned       OrderStatus = "returned"
	OrderStatusExpired        OrderStatus = "expired"
	OrderStatusCancelled      OrderStatus = "cancelled"
)

// TemporalType represents the type of temporal element.
type TemporalType string

const (
	TemporalTypeTrial        TemporalType = "trial"
	TemporalTypeEvaluation   TemporalType = "evaluation"
	TemporalTypeSubscription TemporalType = "subscription"
	TemporalTypeWarranty     TemporalType = "warranty"
	TemporalTypeGuarantee    TemporalType = "guarantee"
)

// TemporalStatus represents the status of a temporal element.
type TemporalStatus string

const (
	TemporalStatusPending       TemporalStatus = "pending"
	TemporalStatusActive        TemporalStatus = "active"
	TemporalStatusExpiring      TemporalStatus = "expiring"
	TemporalStatusExpired       TemporalStatus = "expired"
	TemporalStatusConverted     TemporalStatus = "converted"
	TemporalStatusCancelled     TemporalStatus = "cancelled"
	TemporalStatusSuspended     TemporalStatus = "suspended"
	TemporalStatusRenewed       TemporalStatus = "renewed"
	TemporalStatusClaimed       TemporalStatus = "claimed"
	TemporalStatusClaimApproved TemporalStatus = "claim_approved"
	TemporalStatusClaimDenied   TemporalStatus = "claim_denied"
)

// ShipmentStatus represents the status of a shipment.
type ShipmentStatus string

const (
	ShipmentStatusPending          ShipmentStatus = "pending"
	ShipmentStatusInTransit        ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery   ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered        ShipmentStatus = "delivered"
	ShipmentStatusFailed           ShipmentStatus = "failed"
	ShipmentStatusReturnedToSender ShipmentStatus = "returned_to_sender"
)

// ProductType represents the type of product.
type ProductType string

const (
	ProductTypePhysical ProductType = "physical"
	ProductTypeDigital  ProductType = "digital"
)

// Carrier represents a shipping carrier.
type Carrier string

const (
	CarrierUPS       Carrier = "ups"
	CarrierFedEx     Carrier = "fedex"
	CarrierUSPS      Carrier = "usps"
	CarrierDHL       Carrier = "dhl"
	CarrierShippo    Carrier = "shippo"
	CarrierAfterShip Carrier = "aftership"
)

// DurationUnit represents the unit of a temporal duration.
type DurationUnit string

const (
	DurationUnitDays   DurationUnit = "days"
	DurationUnitWeeks  DurationUnit = "weeks"
	DurationUnitMonths DurationUnit = "months"
	DurationUnitYears  DurationUnit = "years"
)

// DeliverySource represents how a delivery was confirmed.
type DeliverySource string

const (
	DeliverySourceWebhook         DeliverySource = "webhook"
	DeliverySourcePoll            DeliverySource = "poll"
	DeliverySourceManual          DeliverySource = "manual"
	DeliverySourceFallbackCarrier DeliverySource = "fallback_carrier"
)

// WebhookEvent represents a webhook event type.
type WebhookEvent string

const (
	WebhookEventOrderCreated          WebhookEvent = "order.created"
	WebhookEventOrderDelivered        WebhookEvent = "order.delivered"
	WebhookEventDeliveryFailed        WebhookEvent = "delivery.failed"
	WebhookEventTrialStarted          WebhookEvent = "trial.started"
	WebhookEventTrialExpiring         WebhookEvent = "trial.expiring"
	WebhookEventTrialExpired          WebhookEvent = "trial.expired"
	WebhookEventTrialConverted        WebhookEvent = "trial.converted"
	WebhookEventCancellationInitiated WebhookEvent = "cancellation.initiated"
	WebhookEventRiskScoreChanged      WebhookEvent = "risk_score.changed"
	WebhookEventSubscriptionRenewed   WebhookEvent = "subscription.renewed"
	WebhookEventWarrantyClaimed       WebhookEvent = "warranty.claimed"
	WebhookEventTemporalExtended      WebhookEvent = "temporal.extended"
	WebhookEventTemporalAdjusted      WebhookEvent = "temporal.adjusted"
	WebhookEventWarrantyClaimResolved WebhookEvent = "warranty.claim_resolved"
	WebhookEventPaymentSucceeded      WebhookEvent = "payment.succeeded"
	WebhookEventPaymentFailed         WebhookEvent = "payment.failed"
	WebhookEventDisputeCreated        WebhookEvent = "dispute.created"
	WebhookEventDisputeWon            WebhookEvent = "dispute.won"
	WebhookEventDisputeLost           WebhookEvent = "dispute.lost"
)

// ---------------------------------------------------------------------------
// Core Resources
// ---------------------------------------------------------------------------

// Order represents a TrueTrial order.
type Order struct {
	ID               string      `json:"id"`
	TenantID         string      `json:"tenant_id"`
	ConsumerID       string      `json:"consumer_id"`
	ExternalOrderID  string      `json:"external_order_id"`
	ProductName      string      `json:"product_name"`
	ProductType      ProductType `json:"product_type"`
	ProductPriceCents int        `json:"product_price_cents"`
	ProductCurrency  string      `json:"product_currency"`
	Status           OrderStatus `json:"status"`
	Metadata         any         `json:"metadata,omitempty"`
	CreatedAt        string      `json:"created_at"`
	Consumer         *Consumer   `json:"consumer,omitempty"`
}

// Consumer represents a TrueTrial consumer.
type Consumer struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

// Shipment represents a shipment associated with an order.
type Shipment struct {
	ID                string         `json:"id"`
	OrderID           string         `json:"order_id"`
	Carrier           Carrier        `json:"carrier"`
	TrackingNumber    string         `json:"tracking_number"`
	Status            ShipmentStatus `json:"status"`
	EstimatedDelivery string         `json:"estimated_delivery,omitempty"`
	DeliveredAt       string         `json:"delivered_at,omitempty"`
}

// TemporalElement represents a trial, warranty, subscription, evaluation,
// or guarantee period attached to an order.
type TemporalElement struct {
	ID            string         `json:"id"`
	OrderID       string         `json:"order_id"`
	Type          TemporalType   `json:"type"`
	Status        TemporalStatus `json:"status"`
	DurationValue int            `json:"duration_value"`
	DurationUnit  DurationUnit   `json:"duration_unit"`
	BeginTime     string         `json:"begin_time,omitempty"`
	EndTime       string         `json:"end_time,omitempty"`
	ExpiresAt     string         `json:"expires_at,omitempty"`
}

// Cancellation represents a cancellation request for an order.
type Cancellation struct {
	ID          string `json:"id"`
	OrderID     string `json:"order_id"`
	Reason      string `json:"reason"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelled_at,omitempty"`
}

// WebhookSubscription represents a registered webhook endpoint.
type WebhookSubscription struct {
	ID              string         `json:"id"`
	URL             string         `json:"url"`
	Events          []WebhookEvent `json:"events"`
	Secret          string         `json:"secret,omitempty"`
	LastTriggeredAt string         `json:"last_triggered_at,omitempty"`
}

// OrderStatusResponse contains the combined status of an order,
// its temporal element, and its shipment.
type OrderStatusResponse struct {
	OrderID               string         `json:"order_id"`
	ExternalOrderID       string         `json:"external_order_id"`
	OrderStatus           OrderStatus    `json:"order_status"`
	TemporalElementStatus TemporalStatus `json:"temporal_element_status,omitempty"`
	ShipmentStatus        ShipmentStatus `json:"shipment_status,omitempty"`
}

// CarrierHealthEntry represents the health status of a single carrier.
type CarrierHealthEntry struct {
	Carrier   Carrier `json:"carrier"`
	Status    string  `json:"status"`
	Latency   int     `json:"latency_ms"`
	CheckedAt string  `json:"checked_at"`
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

// PaginatedResponse wraps paginated API responses. Data contains the raw JSON
// array of results; use json.Unmarshal to decode into the target slice type.
type PaginatedResponse struct {
	Data        json.RawMessage `json:"data"`
	CurrentPage int             `json:"current_page"`
	LastPage    int             `json:"last_page"`
	PerPage     int             `json:"per_page"`
	Total       int             `json:"total"`
}

// ---------------------------------------------------------------------------
// Request Params
// ---------------------------------------------------------------------------

// ConsumerParams contains the consumer details when creating an order.
type ConsumerParams struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

// CreateOrderParams contains the fields for creating a new order.
type CreateOrderParams struct {
	ExternalOrderID string       `json:"external_order_id"`
	ProductName     string       `json:"product_name"`
	ProductType     ProductType  `json:"product_type"`
	PriceCents      int          `json:"product_price_cents"`
	Currency        string       `json:"product_currency"`
	TemporalType    TemporalType `json:"temporal_type"`
	DurationValue   int          `json:"duration_value"`
	DurationUnit    DurationUnit `json:"duration_unit"`
	Consumer        ConsumerParams `json:"consumer"`
	Metadata        any          `json:"metadata,omitempty"`
}

// ListOrdersParams contains optional query parameters for listing orders.
type ListOrdersParams struct {
	Page   int         `json:"page,omitempty"`
	Status OrderStatus `json:"status,omitempty"`
}

// CreateShipmentParams contains the fields for creating a shipment.
type CreateShipmentParams struct {
	Carrier           Carrier `json:"carrier"`
	TrackingNumber    string  `json:"tracking_number"`
	EstimatedDelivery string  `json:"estimated_delivery,omitempty"`
}

// ConfirmDigitalDeliveryParams contains the fields for confirming
// a digital product delivery.
type ConfirmDigitalDeliveryParams struct {
	DeliveredAt    string         `json:"delivered_at,omitempty"`
	DeliverySource DeliverySource `json:"delivery_source,omitempty"`
}

// ExtendTemporalParams contains the fields for extending a temporal element.
type ExtendTemporalParams struct {
	DurationValue int          `json:"duration_value"`
	DurationUnit  DurationUnit `json:"duration_unit"`
	Reason        string       `json:"reason,omitempty"`
}

// AdjustTemporalParams contains the fields for adjusting a temporal element.
type AdjustTemporalParams struct {
	NewEndTime string `json:"new_end_time"`
	Reason     string `json:"reason,omitempty"`
}

// ClaimParams contains the fields for submitting a warranty or guarantee claim.
type ClaimParams struct {
	Reason      string `json:"reason"`
	Description string `json:"description,omitempty"`
}

// ResolveClaimParams contains the fields for resolving a claim.
type ResolveClaimParams struct {
	Resolution string `json:"resolution"`
	Notes      string `json:"notes,omitempty"`
}

// CreateCancellationParams contains the fields for creating a cancellation.
type CreateCancellationParams struct {
	Reason string `json:"reason"`
}

// CreateWebhookParams contains the fields for creating a webhook subscription.
type CreateWebhookParams struct {
	URL    string         `json:"url"`
	Events []WebhookEvent `json:"events"`
}
