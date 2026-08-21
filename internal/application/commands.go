package application

import (
	"context"
	"example.com/media-workflow/internal/domain"
	"fmt"
	"time"
)

type CommandLog struct {
	ID          string
	Type        string
	AggregateID string
	At          time.Time
}
type CommandBus struct{ log []CommandLog }

func NewCommandBus() *CommandBus { return &CommandBus{log: []CommandLog{}} }

func (b *CommandBus) Recent(limit int) []CommandLog {
	if limit <= 0 || limit >= len(b.log) {
		return b.log
	}
	return b.log[len(b.log)-limit:]
}

func (b *CommandBus) Dispatch(_ context.Context, id, typ, aggregate string) error {
	if id == "" || typ == "" || aggregate == "" {
		return fmt.Errorf("%w: command fields", domain.ErrInvalidInput)
	}
	for _, v := range b.log {
		if v.ID == id {
			return domain.ErrConflict
		}
	}
	b.log = append(b.log, CommandLog{ID: id, Type: typ, AggregateID: aggregate, At: time.Now()})
	return nil
}
func (b *CommandBus) List() []CommandLog {
	return b.log
}
