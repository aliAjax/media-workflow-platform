package security

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"path/filepath"
	"strings"
)

func SafeJoin(root, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("%w: empty path", domain.ErrInvalidInput)
	}
	clean := filepath.Clean(name)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: path traversal", domain.ErrInvalidInput)
	}
	return filepath.Join(root, clean), nil
}
func ValidateContentType(v string) error {
	allowed := map[string]bool{"video/mp4": true, "video/webm": true, "audio/mpeg": true, "audio/wav": true, "application/pdf": true}
	if !allowed[v] {
		return fmt.Errorf("%w: unsupported content type", domain.ErrInvalidInput)
	}
	return nil
}
