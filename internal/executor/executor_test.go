package executor

import (
	"context"
	"errors"
	"example.com/media-workflow/internal/domain"
	"testing"
)

func TestExecutorHonorsCanceledContext(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := New().Run(c, domain.Stage{ID: "s", Kind: "probe", TimeoutSeconds: 1}, domain.Asset{ID: "a"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}
