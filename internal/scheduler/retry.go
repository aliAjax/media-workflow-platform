package scheduler

import (
	"math"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Jitter      bool
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = time.Second
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 5 * time.Minute
	}
	d := float64(p.BaseDelay) * math.Pow(2, float64(attempt-1))
	if d > float64(p.MaxDelay) {
		d = float64(p.MaxDelay)
	}
	if p.Jitter {
		d *= 0.75 + rand.Float64()*0.5
	}
	return time.Duration(d)
}
func (p RetryPolicy) CanRetry(attempt int) bool {
	return p.MaxAttempts <= 0 || attempt <= p.MaxAttempts
}
