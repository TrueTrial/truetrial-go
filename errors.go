package truetrial

import "fmt"

// TrueTrialError is the base error type returned by the TrueTrial API client.
// All API errors can be type-asserted to *TrueTrialError.
type TrueTrialError struct {
	// StatusCode is the HTTP status code from the API response.
	StatusCode int `json:"status_code"`

	// Message is the human-readable error message.
	Message string `json:"message"`

	// ResponseBody is the raw response body for debugging.
	ResponseBody []byte `json:"-"`

	// Errors contains field-level validation errors (only for 422 responses).
	Errors map[string][]string `json:"errors,omitempty"`

	// RetryAfter is the number of seconds to wait before retrying
	// (only for 429 responses).
	RetryAfter int `json:"retry_after,omitempty"`
}

// Error implements the error interface.
func (e *TrueTrialError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("truetrial: %d %s (validation errors: %v)", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("truetrial: %d %s", e.StatusCode, e.Message)
}

// NewAuthenticationError creates an error for 401/403 responses.
func NewAuthenticationError(message string, body []byte) *TrueTrialError {
	return &TrueTrialError{
		StatusCode:   401,
		Message:      message,
		ResponseBody: body,
	}
}

// NewValidationError creates an error for 422 responses with field-level errors.
func NewValidationError(message string, errors map[string][]string, body []byte) *TrueTrialError {
	return &TrueTrialError{
		StatusCode:   422,
		Message:      message,
		Errors:       errors,
		ResponseBody: body,
	}
}

// NewNotFoundError creates an error for 404 responses.
func NewNotFoundError(message string, body []byte) *TrueTrialError {
	return &TrueTrialError{
		StatusCode:   404,
		Message:      message,
		ResponseBody: body,
	}
}

// NewRateLimitError creates an error for 429 responses.
func NewRateLimitError(message string, retryAfter int, body []byte) *TrueTrialError {
	return &TrueTrialError{
		StatusCode:   429,
		Message:      message,
		RetryAfter:   retryAfter,
		ResponseBody: body,
	}
}

// NewServerError creates an error for 5xx responses.
func NewServerError(statusCode int, message string, body []byte) *TrueTrialError {
	return &TrueTrialError{
		StatusCode:   statusCode,
		Message:      message,
		ResponseBody: body,
	}
}

// IsAuthenticationError reports whether the error is an authentication error (401/403).
func IsAuthenticationError(err error) bool {
	e, ok := err.(*TrueTrialError)
	return ok && (e.StatusCode == 401 || e.StatusCode == 403)
}

// IsValidationError reports whether the error is a validation error (422).
func IsValidationError(err error) bool {
	e, ok := err.(*TrueTrialError)
	return ok && e.StatusCode == 422
}

// IsNotFoundError reports whether the error is a not-found error (404).
func IsNotFoundError(err error) bool {
	e, ok := err.(*TrueTrialError)
	return ok && e.StatusCode == 404
}

// IsRateLimitError reports whether the error is a rate-limit error (429).
func IsRateLimitError(err error) bool {
	e, ok := err.(*TrueTrialError)
	return ok && e.StatusCode == 429
}

// IsServerError reports whether the error is a server error (5xx).
func IsServerError(err error) bool {
	e, ok := err.(*TrueTrialError)
	return ok && e.StatusCode >= 500
}
