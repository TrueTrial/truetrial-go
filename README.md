# TrueTrial Go SDK

Official Go client library for the [TrueTrial](https://truetrial.com) API. Manage trial periods, warranties, subscriptions, and guarantees that begin on confirmed delivery -- not order placement.

## Requirements

- Go 1.21 or later
- No external dependencies (stdlib only)

## Installation

```bash
go get github.com/truetrial/truetrial-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/truetrial/truetrial-go"
)

func main() {
    client := truetrial.NewClient("your-api-key")
    ctx := context.Background()

    // Create an order
    order, err := client.Orders.Create(ctx, truetrial.CreateOrderParams{
        ExternalOrderID: "EXT-12345",
        ProductName:     "Premium Supplement",
        ProductType:     truetrial.ProductTypePhysical,
        PriceCents:      4999,
        Currency:        "USD",
        TemporalType:    truetrial.TemporalTypeTrial,
        DurationValue:   30,
        DurationUnit:    truetrial.DurationUnitDays,
        Consumer: truetrial.ConsumerParams{
            Email:     "customer@example.com",
            FirstName: "Jane",
            LastName:  "Doe",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created order: %s\n", order.ID)
}
```

## Configuration

```go
// Default configuration
client := truetrial.NewClient("your-api-key")

// Custom base URL
client := truetrial.NewClient("your-api-key",
    truetrial.WithBaseURL("https://api.truetrial.com/api/v1"),
)

// Custom HTTP client (e.g., for proxies or custom TLS)
client := truetrial.NewClient("your-api-key",
    truetrial.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)
```

## API Reference

All methods accept a `context.Context` as their first parameter and return an error as their last return value.

### Orders

```go
// List orders with optional filters
resp, err := client.Orders.List(ctx, &truetrial.ListOrdersParams{
    Page:   1,
    Status: truetrial.OrderStatusDelivered,
})

// Create an order
order, err := client.Orders.Create(ctx, truetrial.CreateOrderParams{...})

// Get a single order
order, err := client.Orders.Get(ctx, "order-id")

// Get combined status (order + temporal + shipment)
status, err := client.Orders.Status(ctx, "order-id")
```

### Shipments

```go
// Create a shipment
shipment, err := client.Shipments.Create(ctx, "order-id", truetrial.CreateShipmentParams{
    Carrier:        truetrial.CarrierUPS,
    TrackingNumber: "1Z999AA10123456784",
})

// List shipments for an order
shipments, err := client.Shipments.List(ctx, "order-id")
```

### Digital Delivery

```go
// Confirm digital product delivery
order, err := client.DigitalDelivery.Confirm(ctx, "order-id", truetrial.ConfirmDigitalDeliveryParams{
    DeliverySource: truetrial.DeliverySourceWebhook,
})
```

### Temporal Elements

```go
// Get temporal element for an order
temporal, err := client.Temporal.Get(ctx, "order-id")

// Extend a trial/warranty/etc.
temporal, err := client.Temporal.Extend(ctx, "order-id", truetrial.ExtendTemporalParams{
    DurationValue: 7,
    DurationUnit:  truetrial.DurationUnitDays,
    Reason:        "Customer requested extension",
})

// Adjust the end time directly
temporal, err := client.Temporal.Adjust(ctx, "order-id", truetrial.AdjustTemporalParams{
    NewEndTime: "2026-03-15T00:00:00Z",
    Reason:     "Shipping delay compensation",
})

// Submit a warranty/guarantee claim
temporal, err := client.Temporal.Claim(ctx, "order-id", truetrial.ClaimParams{
    Reason:      "Product defective",
    Description: "Screen cracked on arrival",
})

// Resolve a claim
temporal, err := client.Temporal.ResolveClaim(ctx, "order-id", truetrial.ResolveClaimParams{
    Resolution: "approved",
    Notes:      "Replacement shipped",
})
```

### Cancellations

```go
// Create a cancellation
cancellation, err := client.Cancellations.Create(ctx, "order-id", truetrial.CreateCancellationParams{
    Reason: "Customer requested cancellation",
})

// Get cancellation details
cancellation, err := client.Cancellations.Get(ctx, "order-id")
```

### Webhooks

```go
// List webhook subscriptions
webhooks, err := client.Webhooks.List(ctx)

// Create a webhook subscription
webhook, err := client.Webhooks.Create(ctx, truetrial.CreateWebhookParams{
    URL: "https://example.com/webhooks/truetrial",
    Events: []truetrial.WebhookEvent{
        truetrial.WebhookEventOrderDelivered,
        truetrial.WebhookEventTrialExpiring,
        truetrial.WebhookEventTrialExpired,
    },
})
// Store webhook.Secret securely for signature verification.

// Delete a webhook subscription
err := client.Webhooks.Delete(ctx, "webhook-id")
```

### System

```go
// Check carrier health
carriers, err := client.System.CarrierHealth(ctx)
for _, c := range carriers {
    fmt.Printf("%s: %s (%dms)\n", c.Carrier, c.Status, c.Latency)
}
```

## Error Handling

All API errors are returned as `*truetrial.TrueTrialError` with typed helper functions for common cases:

```go
order, err := client.Orders.Get(ctx, "order-id")
if err != nil {
    if truetrial.IsNotFoundError(err) {
        // Order does not exist
        log.Println("Order not found")
        return
    }

    if truetrial.IsValidationError(err) {
        // Access field-level validation errors
        ttErr := err.(*truetrial.TrueTrialError)
        for field, messages := range ttErr.Errors {
            log.Printf("  %s: %v\n", field, messages)
        }
        return
    }

    if truetrial.IsRateLimitError(err) {
        ttErr := err.(*truetrial.TrueTrialError)
        log.Printf("Rate limited. Retry after %d seconds\n", ttErr.RetryAfter)
        return
    }

    if truetrial.IsAuthenticationError(err) {
        log.Fatal("Invalid API key")
    }

    if truetrial.IsServerError(err) {
        log.Println("TrueTrial server error, try again later")
        return
    }

    log.Fatal(err)
}
```

## Webhook Verification

Verify incoming webhook signatures to ensure they are authentic:

```go
import "github.com/truetrial/truetrial-go"

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    payload, _ := io.ReadAll(r.Body)
    signature := r.Header.Get(truetrial.WebhookSignatureHeader)
    event := r.Header.Get(truetrial.WebhookEventHeader)
    secret := "your-webhook-secret"

    // Simple signature verification
    if !truetrial.VerifyWebhookSignature(payload, signature, secret) {
        http.Error(w, "Invalid signature", http.StatusForbidden)
        return
    }

    // Or verify with timestamp tolerance (prevents replay attacks)
    if !truetrial.VerifyWebhookSignatureWithTolerance(payload, signature, secret, 5*time.Minute) {
        http.Error(w, "Invalid or expired signature", http.StatusForbidden)
        return
    }

    switch truetrial.WebhookEvent(event) {
    case truetrial.WebhookEventOrderDelivered:
        // Handle delivery confirmation
    case truetrial.WebhookEventTrialExpiring:
        // Send reminder to customer
    case truetrial.WebhookEventTrialExpired:
        // Process conversion or return
    }

    w.WriteHeader(http.StatusOK)
}
```

## Pagination

List endpoints return a `PaginatedResponse` with raw JSON data that you decode into the target type:

```go
import "encoding/json"

resp, err := client.Orders.List(ctx, &truetrial.ListOrdersParams{Page: 1})
if err != nil {
    log.Fatal(err)
}

var orders []truetrial.Order
if err := json.Unmarshal(resp.Data, &orders); err != nil {
    log.Fatal(err)
}

fmt.Printf("Page %d of %d (%d total orders)\n", resp.CurrentPage, resp.LastPage, resp.Total)
```

## License

See [LICENSE](LICENSE) for details.
