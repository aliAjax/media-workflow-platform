package observability

import (
	"context"
	"time"
)

type Check struct {
	Name string
	Run  func(context.Context) error
}
type Report struct {
	Status   string
	Checks   map[string]string
	Duration time.Duration
}

func RunChecks(ctx context.Context, checks []Check) Report {
	start := time.Now()
	r := Report{Status: "ok", Checks: map[string]string{}}
	for _, c := range checks {
		e := c.Run(ctx)
		if e != nil {
			r.Status = "degraded"
			r.Checks[c.Name] = e.Error()
		} else {
			r.Checks[c.Name] = "ok"
		}
	}
	r.Duration = time.Since(start)
	return r
}
