package domain

import (
	"fmt"
	"sync"
	"time"
)

type ProgressEvent struct {
	JobID   string
	Stage   string
	Percent int
	Message string
	At      time.Time
}
type ProgressTracker struct {
	mu     sync.RWMutex
	events map[string][]ProgressEvent
	max    int
}

func NewProgressTracker(max int) *ProgressTracker {
	if max < 1 {
		max = 100
	}
	return &ProgressTracker{events: map[string][]ProgressEvent{}, max: max}
}
func (t *ProgressTracker) Publish(e ProgressEvent) error {
	if e.JobID == "" || e.Percent < 0 || e.Percent > 100 {
		return fmt.Errorf("%w: invalid progress", ErrInvalidInput)
	}
	if e.At.IsZero() {
		e.At = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	list := t.events[e.JobID]
	if len(list) >= t.max {
		list = list[1:]
	}
	t.events[e.JobID] = append(list, e)
	return nil
}
func (t *ProgressTracker) List(job string) []ProgressEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	src := t.events[job]
	out := make([]ProgressEvent, len(src))
	copy(out, src)
	return out
}
