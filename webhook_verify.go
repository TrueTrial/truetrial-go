package truetrial

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// WebhookSignatureHeader is the HTTP header containing the webhook signature.
	WebhookSignatureHeader = "X-TrueTrial-Signature"

	// WebhookEventHeader is the HTTP header containing the webhook event type.
	WebhookEventHeader = "X-TrueTrial-Event"

	// WebhookTimestampHeader is the HTTP header containing the webhook timestamp.
	WebhookTimestampHeader = "X-TrueTrial-Timestamp"

	// DefaultWebhookTolerance is the default maximum age allowed for a webhook
	// delivery before it is considered stale. Five minutes is the default.
	DefaultWebhookTolerance = 5 * time.Minute
)

// VerifyWebhookSignature verifies that a webhook payload was signed by
// TrueTrial using the given signing secret. The signature should be the
// value of the X-TrueTrial-Signature header.
//
// This performs a constant-time comparison to prevent timing attacks.
func VerifyWebhookSignature(payload []byte, signature string, secret string) bool {
	expected := computeHMAC(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expected))
}

// VerifyWebhookSignatureWithTolerance verifies the webhook signature and
// additionally checks that the timestamp in the signature is within the
// given tolerance window.
//
// The signature format is expected to be "t=<unix_timestamp>,v1=<hex_hmac>"
// where the HMAC is computed over "<timestamp>.<payload>".
//
// If tolerance is zero, DefaultWebhookTolerance is used.
func VerifyWebhookSignatureWithTolerance(payload []byte, signature string, secret string, tolerance time.Duration) bool {
	if tolerance == 0 {
		tolerance = DefaultWebhookTolerance
	}

	parts := parseSignature(signature)
	timestamp := parts["t"]
	hash := parts["v1"]

	if timestamp == "" || hash == "" {
		return false
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	// Check timestamp is within tolerance.
	age := math.Abs(float64(time.Now().Unix() - ts))
	if age > tolerance.Seconds() {
		return false
	}

	// Compute expected HMAC over "timestamp.payload".
	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(payload))
	expected := computeHMAC([]byte(signedPayload), secret)

	return hmac.Equal([]byte(hash), []byte(expected))
}

// computeHMAC computes the HMAC-SHA256 of the data using the given key
// and returns the result as a lowercase hex string.
func computeHMAC(data []byte, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// parseSignature splits a structured signature string like
// "t=123456,v1=abcdef" into its components.
func parseSignature(sig string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(sig, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}
