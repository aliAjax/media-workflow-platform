package pipeline

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"strconv"
	"strings"
)

func ParseInt(params map[string]string, key string, defaultValue int) (int, error) {
	v := strings.TrimSpace(params[key])
	if v == "" {
		return defaultValue, nil
	}
	n, e := strconv.Atoi(v)
	if e != nil {
		return 0, fmt.Errorf("%s: %v", domain.ErrInvalidInput, e)
	}
	return n, nil
}
func ParseBool(params map[string]string, key string, defaultValue bool) (bool, error) {
	v := strings.TrimSpace(params[key])
	if v == "" {
		return defaultValue, nil
	}
	n, e := strconv.ParseBool(v)
	if e != nil {
		return false, fmt.Errorf("%s: %v", domain.ErrInvalidInput, e)
	}
	return n, nil
}
func AllowlistedParams(params map[string]string, allowed map[string]bool) error {
	for k := range params {
		if !allowed[k] {
			return fmt.Errorf("%s: parameter %s not allowed", domain.ErrInvalidInput, k)
		}
	}
	return nil
}
