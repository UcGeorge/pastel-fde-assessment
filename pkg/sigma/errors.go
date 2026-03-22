package sigma

import (
	"errors"
	"fmt"
)

// Define Sentinel Errors for Downstream Clients
var (
	ErrBadRequest       = errors.New("bad request")
	ErrUnauthorized     = errors.New("unauthorized request")
	ErrResourceNotFound = errors.New("request resource does not exist")
	ErrRateLimited      = errors.New("too many requests sent to the API")
	ErrInternalServer   = errors.New("sigma internal server error")
	ErrUnknown          = errors.New("unknown api error")
)

// Define the Custom Error Struct
type SigmaAPIError struct {
	StatusCode int    `json:"-"`       // HTTP Status Code
	ErrorCode  any    `json:"error"`   // e.g., 123 or "123"
	Message    string `json:"message"` // The human-readable error from Sigma
	Details    any    `json:"details"` // Optional map or slice of specific field errors

	// The internal Sentinel error to wrap (allows errors.Is to work)
	err error
}

// Implement the `error` interface
func (e *SigmaAPIError) Error() string {
	if e.ErrorCode != nil {
		return fmt.Sprintf("sigma api error [%d %v]: %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("sigma api error [%d]: %s", e.StatusCode, e.Message)
}

// Implement the `Unwrap` interface
func (e *SigmaAPIError) Unwrap() error {
	return e.err
}
