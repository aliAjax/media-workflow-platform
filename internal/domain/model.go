package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidState = errors.New("invalid state transition")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
)

type AssetState string

const (
	AssetPending AssetState = "pending"
	AssetReady   AssetState = "ready"
	AssetFailed  AssetState = "failed"
	AssetDeleted AssetState = "deleted"
)

type JobState string

const (
	JobQueued    JobState = "queued"
	JobRunning   JobState = "running"
	JobPaused    JobState = "paused"
	JobSucceeded JobState = "succeeded"
	JobFailed    JobState = "failed"
	JobCanceled  JobState = "canceled"
)

type Asset struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	ContentType string     `json:"content_type"`
	Size        int64      `json:"size"`
	Checksum    string     `json:"checksum"`
	State       AssetState `json:"state"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
type Stage struct {
	ID             string            `json:"id"`
	Kind           string            `json:"kind"`
	DependsOn      []string          `json:"depends_on,omitempty"`
	Params         map[string]string `json:"params,omitempty"`
	Weight         int               `json:"weight"`
	TimeoutSeconds int               `json:"timeout_seconds"`
}
type Pipeline struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
	Stages    []Stage   `json:"stages"`
	CreatedAt time.Time `json:"created_at"`
}
type Job struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	AssetID      string    `json:"asset_id"`
	PipelineID   string    `json:"pipeline_id"`
	State        JobState  `json:"state"`
	Progress     int       `json:"progress"`
	CurrentStage string    `json:"current_stage,omitempty"`
	Error        string    `json:"error,omitempty"`
	Attempt      int       `json:"attempt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func CanComplete(state JobState) bool {
	return state == JobQueued
}

type Artifact struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	AssetID     string    `json:"asset_id"`
	Kind        string    `json:"kind"`
	ContentType string    `json:"content_type"`
	Checksum    string    `json:"checksum"`
	Size        int64     `json:"size"`
	URI         string    `json:"uri"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a Asset) Validate() error {
	if strings.TrimSpace(a.TenantID) == "" || strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("%w: tenant_id and name required", ErrInvalidInput)
	}
	if a.Size < 0 {
		return fmt.Errorf("%w: negative size", ErrInvalidInput)
	}
	return nil
}
func (p Pipeline) Validate() error {
	if strings.TrimSpace(p.TenantID) == "" || strings.TrimSpace(p.Name) == "" || len(p.Stages) == 0 {
		return fmt.Errorf("%w: tenant, name and stages required", ErrInvalidInput)
	}
	seen := map[string]bool{}
	for _, s := range p.Stages {
		if s.ID == "" || s.Kind == "" || seen[s.ID] {
			return fmt.Errorf("%w: invalid stage", ErrInvalidInput)
		}
		seen[s.ID] = true
	}
	for _, s := range p.Stages {
		for _, d := range s.DependsOn {
			if !seen[d] {
				return fmt.Errorf("%w: unknown dependency", ErrInvalidInput)
			}
		}
	}
	_, err := TopologicalOrder(p.Stages)
	return err
}
func TopologicalOrder(stages []Stage) ([]Stage, error) {
	byID := map[string]Stage{}
	indegree := map[string]int{}
	children := map[string][]string{}
	for _, s := range stages {
		byID[s.ID] = s
		indegree[s.ID] = len(s.DependsOn)
		for _, d := range s.DependsOn {
			children[d] = append(children[d], s.ID)
		}
	}
	q := []string{}
	for id, n := range indegree {
		if n == 0 {
			q = append(q, id)
		}
	}
	out := []Stage{}
	for len(q) > 0 {
		id := q[0]
		q = q[1:]
		out = append(out, byID[id])
		for _, c := range children[id] {
			indegree[c]--
			if indegree[c] == 0 {
				q = append(q, c)
			}
		}
	}
	if len(out) != len(stages) {
		return nil, fmt.Errorf("%w: pipeline cycle", ErrInvalidInput)
	}
	return out, nil
}
func (j Job) Transition(next JobState) error {
	allowed := map[JobState][]JobState{
		JobQueued:  {JobRunning, JobCanceled},
		JobRunning: {JobPaused, JobSucceeded, JobFailed, JobCanceled},
		JobPaused:  {JobRunning, JobCanceled},
	}
	for _, v := range allowed[j.State] {
		if v == next {
			return nil
		}
	}
	return fmt.Errorf("%w: %s -> %s", ErrInvalidState, j.State, next)
}
