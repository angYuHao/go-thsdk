package thsdk

import (
	"errors"
	"fmt"
)

// ErrNotAuthenticated is returned by data calls before Auth or AuthQrcode succeeds.
var ErrNotAuthenticated = errors.New("thsdk: 尚未认证；请先调用 Auth 或 AuthQrcode")

// Error is a business error reported by the runtime component.
type Error struct {
	Code          string
	Message       string
	RetryAfter    int // milliseconds, when the component throttled the call
	HasRetryAfter bool
}

func (e *Error) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return e.Code + ": " + e.Message
}

// IsRateLimited reports whether the component asked for a slower call rate.
func (e *Error) IsRateLimited() bool { return e != nil && e.Code == "rate_limited" }

// AuthError is any authentication or session-restore failure.
type AuthError struct {
	Reason string
	Err    error
}

func (e *AuthError) Error() string {
	if e.Err == nil {
		return "thsdk: " + e.Reason
	}
	return fmt.Sprintf("thsdk: %s: %v", e.Reason, e.Err)
}

func (e *AuthError) Unwrap() error { return e.Err }

// ParameterError is a client-side request validation failure.
type ParameterError struct {
	Argument string
	Message  string
	Err      error
}

func (e *ParameterError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("thsdk: %s %s: %v", e.Argument, e.Message, e.Err)
	}
	return fmt.Sprintf("thsdk: %s %s", e.Argument, e.Message)
}

func (e *ParameterError) Unwrap() error { return e.Err }

func invalidArgument(argument, format string, args ...any) error {
	return &ParameterError{Argument: argument, Message: fmt.Sprintf(format, args...)}
}

func newAuthError(reason string, err error) error {
	return &AuthError{Reason: reason, Err: err}
}
