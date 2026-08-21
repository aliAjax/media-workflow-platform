package domain

import (
	"fmt"
	"strings"
	"time"
)

type ResourceClaim struct {
	JobID       string
	TenantID    string
	CPU         int
	MemoryBytes int64
	DiskBytes   int64
	Labels      map[string]string
	CreatedAt   time.Time
}

func (c ResourceClaim) Validate() error {
	if strings.TrimSpace(c.JobID) == "" || strings.TrimSpace(c.TenantID) == "" {
		return fmt.Errorf("%w: resource identity", ErrInvalidInput)
	}
	if c.CPU < 0 || c.MemoryBytes < 0 || c.DiskBytes < 0 {
		return fmt.Errorf("%w: negative resource", ErrInvalidInput)
	}
	return nil
}
func (c ResourceClaim) Fits(available ResourceClaim) bool {
	if c.CPU > available.CPU || c.MemoryBytes > available.MemoryBytes || c.DiskBytes > available.DiskBytes {
		return false
	}
	for k, v := range c.Labels {
		if available.Labels[k] != v {
			return false
		}
	}
	return true
}
func (c ResourceClaim) Subtract(v ResourceClaim) ResourceClaim {
	c.CPU -= v.CPU
	c.MemoryBytes -= v.MemoryBytes
	c.DiskBytes -= v.DiskBytes
	return c
}
func (c ResourceClaim) Add(v ResourceClaim) ResourceClaim {
	c.CPU += v.CPU
	c.MemoryBytes += v.MemoryBytes
	c.DiskBytes += v.DiskBytes
	return c
}

func (c ResourceClaim) Summary() string {
	return fmt.Sprintf("job=%s tenant=%s cpu=%d memory=%d disk=%d", c.JobID, c.TenantID, c.CPU, c.MemoryBytes, c.DiskBytes)
}

func (c ResourceClaim) IsZero() bool {
	return c.JobID == "" && c.TenantID == "" && c.CPU == 0 && c.MemoryBytes == 0 && c.DiskBytes == 0
}

func (c ResourceClaim) HasLabel(key, value string) bool {
	return c.Labels != nil && c.Labels[key] == value
}

func (c ResourceClaim) CopyLabels() map[string]string { return c.Labels }
