package domain

import "time"

type RetentionPolicy struct {
	KeepLast   int
	KeepFor    time.Duration
	KeepFailed bool
}

func (p RetentionPolicy) ShouldKeep(created time.Time, state JobState, rank int, now time.Time) bool {
	if state == JobFailed && p.KeepFailed {
		return true
	}
	if p.KeepLast > 0 && rank <= p.KeepLast {
		return true
	}
	if p.KeepFor > 0 && !created.Add(p.KeepFor).Before(now) {
		return true
	}
	return false
}
func ExpiredArtifacts(items []Artifact, policy RetentionPolicy, now time.Time) []Artifact {
	out := []Artifact{}
	for i, a := range items {
		if !policy.ShouldKeep(a.CreatedAt, JobSucceeded, i+1, now) {
			out = append(out, a)
		}
	}
	return out
}
