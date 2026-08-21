package domain

import (
	"fmt"
	"sync"
	"time"
)

type Quota struct {
	TenantID      string
	MaxBytes      int64
	UsedBytes     int64
	MaxConcurrent int
	Running       int
	UpdatedAt     time.Time
}
type QuotaLedger struct {
	mu    sync.Mutex
	items map[string]Quota
}

func NewQuotaLedger() *QuotaLedger { return &QuotaLedger{} }
func (q *QuotaLedger) Set(v Quota) error {
	if v.TenantID == "" || v.MaxBytes < 0 || v.MaxConcurrent < 0 {
		return fmt.Errorf("%w: invalid quota", ErrInvalidInput)
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	v.UpdatedAt = time.Now()
	q.items[v.TenantID] = v
	return nil
}
func (q *QuotaLedger) Reserve(t string, bytes int64) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	v, ok := q.items[t]
	if !ok {
		return nil
	}
	if bytes < 0 || v.UsedBytes+bytes > v.MaxBytes {
		return fmt.Errorf("%w: bytes", ErrConflict)
	}
	if v.MaxConcurrent > 0 && v.Running+1 > v.MaxConcurrent {
		return fmt.Errorf("%w: concurrent", ErrConflict)
	}
	v.UsedBytes += bytes
	v.Running++
	v.UpdatedAt = time.Now()
	q.items[t] = v
	return nil
}
func (q *QuotaLedger) Release(t string, bytes int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	v := q.items[t]
	v.UsedBytes -= bytes
	if v.UsedBytes < 0 {
		v.UsedBytes = 0
	}
	if v.Running > 0 {
		v.Running--
	}
	v.UpdatedAt = time.Now()
	q.items[t] = v
}
func (q *QuotaLedger) Get(t string) (Quota, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	v, ok := q.items[t]
	return v, ok
}
