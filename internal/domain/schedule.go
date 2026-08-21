package domain

import (
	"fmt"
	"sort"
	"time"
)

type Schedule struct {
	ID        string
	JobID     string
	RunAt     time.Time
	Timezone  string
	Canceled  bool
	CreatedAt time.Time
}

func (s Schedule) Validate() error {
	if s.ID == "" || s.JobID == "" {
		return fmt.Errorf("%w: schedule identifiers", ErrInvalidInput)
	}
	if s.RunAt.IsZero() {
		return fmt.Errorf("%w: run_at required", ErrInvalidInput)
	}
	return nil
}
func SortSchedules(items []Schedule) {
	sort.Slice(items, func(i, j int) bool { return items[i].RunAt.Before(items[j].RunAt) })
}
func DueSchedules(items []Schedule, now time.Time) []Schedule {
	out := []Schedule{}
	for _, s := range items {
		if !s.Canceled && !s.RunAt.After(now) {
			out = append(out, s)
		}
	}
	return out
}
func NextRun(start time.Time, interval time.Duration, now time.Time) time.Time {
	if interval <= 0 {
		return start
	}
	next := start
	for !next.After(now) {
		next = next.Add(interval)
	}
	return next
}
