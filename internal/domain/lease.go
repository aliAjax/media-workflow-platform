package domain

import (
	"fmt"
	"sync"
	"time"
)

type Lease struct {
	ID         string
	Owner      string
	ExpiresAt  time.Time
	Generation uint64
}
type LeaseManager struct {
	mu     sync.Mutex
	leases map[string]Lease
	ttl    time.Duration
}

func NewLeaseManager(ttl time.Duration) *LeaseManager {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &LeaseManager{leases: map[string]Lease{}, ttl: ttl}
}
func (l *LeaseManager) Acquire(id, owner string) (Lease, error) {
	if id == "" || owner == "" {
		return Lease{}, fmt.Errorf("%w: lease fields", ErrInvalidInput)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	old, ok := l.leases[id]
	if ok && old.ExpiresAt.After(time.Now()) && old.Owner != owner {
		return Lease{}, fmt.Errorf("%w: lease held", ErrConflict)
	}
	v := Lease{ID: id, Owner: owner, ExpiresAt: time.Now().Add(l.ttl), Generation: old.Generation + 1}
	l.leases[id] = v
	return v, nil
}
func (l *LeaseManager) Renew(id, owner string, generation uint64) (Lease, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.leases[id]
	if !ok || v.Owner != owner || v.Generation != generation {
		return Lease{}, fmt.Errorf("%w: stale lease", ErrConflict)
	}
	v.ExpiresAt = time.Now().Add(l.ttl)
	l.leases[id] = v
	return v, nil
}
func (l *LeaseManager) Release(id, owner string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.leases[id]; ok && v.Owner == owner {
		delete(l.leases, id)
	}
}
