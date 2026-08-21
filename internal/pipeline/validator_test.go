package pipeline

import (
	"errors"
	"example.com/media-workflow/internal/domain"
	"testing"
)

func TestParseIntPreservesInvalidInput(t *testing.T) {
	_, err := ParseInt(map[string]string{"width": "bad"}, "width", 0)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error chain lost: %v", err)
	}
}

func TestParseBoolPreservesInvalidInput(t *testing.T) {
	_, err := ParseBool(map[string]string{"enabled": "sometimes"}, "enabled", false)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("boolean error chain lost: %v", err)
	}
}

func TestAllowlistPreservesInvalidInput(t *testing.T) {
	err := AllowlistedParams(map[string]string{"secret": "value"}, map[string]bool{"width": true})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("allowlist error chain lost: %v", err)
	}
}

func TestValidatorPreservesNestedInvalidInput(t *testing.T) {
	if err := ValidateParameter(map[string]string{"width": "bad"}, "width"); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("validator error=%v", err)
	}
}
