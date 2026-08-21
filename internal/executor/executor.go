package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"example.com/media-workflow/internal/domain"
)

type Result struct {
	Kind        string
	ContentType string
	Data        []byte
	Metadata    map[string]string
}
type Executor struct{}

func New() *Executor { return &Executor{} }
func (*Executor) Run(ctx context.Context, stage domain.Stage, asset domain.Asset) (Result, error) {
	d := time.Duration(stage.TimeoutSeconds) * time.Second
	if d <= 0 {
		d = 30 * time.Second
	}
	child, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	select {
	case <-child.Done():
		return Result{}, child.Err()
	case <-time.After(time.Millisecond):
	}
	payload := []byte(fmt.Sprintf("stage=%s asset=%s checksum=%s", stage.Kind, asset.ID, asset.Checksum))
	sum := sha256.Sum256(payload)
	return Result{Kind: stage.Kind, ContentType: "application/octet-stream", Data: payload, Metadata: map[string]string{"sha256": hex.EncodeToString(sum[:])}}, nil
}
