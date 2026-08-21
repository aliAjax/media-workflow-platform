package domain

import (
	"sort"
	"strings"
)

type JobFilter struct {
	TenantID   string
	States     []JobState
	AssetID    string
	PipelineID string
}

func MatchJob(j Job, f JobFilter) bool {
	if f.TenantID != "" && j.TenantID != f.TenantID {
		return false
	}
	if f.AssetID != "" && j.AssetID != f.AssetID {
		return false
	}
	if f.PipelineID != "" && j.PipelineID != f.PipelineID {
		return false
	}
	if len(f.States) > 0 {
		ok := false
		for _, s := range f.States {
			if j.State == s {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	return true
}
func SortJobs(items []Job, field string, desc bool) {
	sort.SliceStable(items, func(i, j int) bool {
		var less bool
		switch strings.ToLower(field) {
		case "progress":
			less = items[i].Progress < items[j].Progress
		case "state":
			less = items[i].State < items[j].State
		default:
			less = items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		if desc {
			return !less
		}
		return less
	})
}
