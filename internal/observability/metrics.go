package observability

import (
	"fmt"
	"sync/atomic"
)

type Counters struct {
	Requests atomic.Uint64
	Jobs     atomic.Uint64
	Failures atomic.Uint64
	Bytes    atomic.Int64
}

func (c *Counters) Snapshot() map[string]any {
	return map[string]any{"requests": c.Requests.Load(), "jobs": c.Jobs.Load(), "failures": c.Failures.Load(), "bytes": c.Bytes.Load()}
}
func FormatSnapshot(v map[string]any) string {
	return fmt.Sprintf("requests=%v jobs=%v failures=%v bytes=%v", v["requests"], v["jobs"], v["failures"], v["bytes"])
}
