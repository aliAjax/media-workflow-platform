package domain

import "fmt"

type ErrorCode string

const (
	CodeInvalidInput ErrorCode = "invalid_input"
	CodeNotFound     ErrorCode = "not_found"
	CodeConflict     ErrorCode = "conflict"
	CodeUnauthorized ErrorCode = "unauthorized"
	CodeUnavailable  ErrorCode = "unavailable"
)

type PublicError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e PublicError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
func (e PublicError) Unwrap() error { return nil }
func NewInvalid(message string, cause error) error {
	return PublicError{Code: CodeInvalidInput, Message: message}
}
func NewNotFound(message string) error { return PublicError{Code: CodeNotFound, Message: message} }
func NewConflict(message string) error { return PublicError{Code: CodeConflict, Message: message} }
func IsRetryable(err error) bool {
	switch err.(type) {
	case PublicError:
		return false
	default:
		return err != nil
	}
}
