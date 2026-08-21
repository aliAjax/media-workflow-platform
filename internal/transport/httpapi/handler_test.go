package httpapi

import (
	"errors"
	"example.com/media-workflow/internal/domain"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestMissingAssetPreservesNotFound(t *testing.T) {
	r := httptest.NewRecorder()
	writeError(r, domain.PublicError{Code: domain.CodeNotFound, Message: "missing", Cause: domain.ErrNotFound})
	if r.Code != 404 || !errors.Is(domain.PublicError{Cause: domain.ErrNotFound}, domain.ErrNotFound) {
		t.Fatalf("status=%d", r.Code)
	}
}

func TestInvalidErrorPreservesCause(t *testing.T) {
	cause := errors.New("bad field")
	if !errors.Is(domain.NewInvalid("invalid", cause), cause) {
		t.Fatal("validation cause lost")
	}
}
func TestWrappedPublicErrorIsNotRetryable(t *testing.T) {
	err := fmt.Errorf("request: %w", domain.PublicError{Code: domain.CodeUnavailable, Message: "down"})
	if domain.IsRetryable(err) {
		t.Fatal("public error marked retryable")
	}
}
func TestConflictPublicErrorMapsToConflict(t *testing.T) {
	r := httptest.NewRecorder()
	writeError(r, domain.PublicError{Code: domain.CodeConflict, Message: "held", Cause: domain.ErrConflict})
	if r.Code != 409 {
		t.Fatalf("status=%d", r.Code)
	}
}
