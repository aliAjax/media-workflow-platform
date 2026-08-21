package pipeline

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"strings"
)

var allowed = map[string]bool{"probe": true, "transcode": true, "thumbnail": true, "waveform": true, "subtitle": true, "package": true, "quality-check": true}

func ValidateParameter(params map[string]string, key string) error {
	_, err := ParseInt(params, key, 0)
	if err != nil {
		return fmt.Errorf("validate %s: %w", key, err)
	}
	return nil
}

func Validate(p domain.Pipeline) error {
	if e := p.Validate(); e != nil {
		return e
	}
	for _, s := range p.Stages {
		if !allowed[s.Kind] {
			return fmt.Errorf("%w: unsupported stage", domain.ErrInvalidInput)
		}
		if s.Weight < 0 || s.Weight > 1000 {
			return fmt.Errorf("%w: weight", domain.ErrInvalidInput)
		}
		for k, v := range s.Params {
			if strings.ContainsAny(k, "\r\n") || strings.ContainsRune(v, 0) {
				return fmt.Errorf("%w: unsafe parameter", domain.ErrInvalidInput)
			}
		}
	}
	return nil
}
