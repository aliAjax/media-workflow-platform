package application

import (
	"encoding/base64"
	"example.com/media-workflow/internal/domain"
	"fmt"
	"strconv"
	"strings"
)

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	Total      int    `json:"total"`
}

func EncodeCursor(offset int) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}
func DecodeCursor(v string) (int, error) {
	if v == "" {
		return 0, nil
	}
	b, e := base64.RawURLEncoding.DecodeString(v)
	if e != nil {
		return 0, e
	}
	n, e := strconv.Atoi(string(b))
	if e != nil {
		return 0, e
	}
	return n, nil
}
func PaginateJobs(items []domain.Job, cursor string, limit int) Page[domain.Job] {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := DecodeCursor(strings.TrimSpace(cursor))
	if offset > len(items) {
		offset = len(items)
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	p := Page[domain.Job]{Items: items[offset:end], Total: len(items)}
	if end < len(items) {
		p.NextCursor = EncodeCursor(end)
	}
	return p
}
func ValidatePage(limit int) error {
	if limit < 0 || limit > 1000 {
		return fmt.Errorf("%w: page limit", domain.ErrInvalidInput)
	}
	return nil
}
