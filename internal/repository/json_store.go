package repository

import (
	"context"
	"encoding/json"
	"example.com/media-workflow/internal/domain"
	"os"
	"path/filepath"
	"sync"
)

type State struct {
	Assets    map[string]domain.Asset    `json:"assets"`
	Pipelines map[string]domain.Pipeline `json:"pipelines"`
	Jobs      map[string]domain.Job      `json:"jobs"`
	Artifacts map[string]domain.Artifact `json:"artifacts"`
}
type Store struct {
	mu    sync.RWMutex
	path  string
	state State
}

func New(path string) (*Store, error) {
	s := &Store{path: path, state: State{map[string]domain.Asset{}, map[string]domain.Pipeline{}, map[string]domain.Job{}, map[string]domain.Artifact{}}}
	if b, e := os.ReadFile(path); e == nil {
		if e = json.Unmarshal(b, &s.state); e != nil {
			return nil, e
		}
	} else if !os.IsNotExist(e) {
		return nil, e
	}
	return s, nil
}
func (s *Store) persist() error {
	b, e := json.MarshalIndent(s.state, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(s.path), 0755); e != nil {
		return e
	}
	tmp := s.path + ".tmp"
	if e = os.WriteFile(tmp, b, 0644); e != nil {
		return e
	}
	return os.Rename(tmp, s.path)
}
func (s *Store) CreateAsset(_ context.Context, a domain.Asset) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Assets[a.ID]; ok {
		return domain.ErrConflict
	}
	s.state.Assets[a.ID] = a
	return s.persist()
}
func (s *Store) GetAsset(_ context.Context, id string) (domain.Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.state.Assets[id]
	if !ok {
		return domain.Asset{}, domain.ErrNotFound
	}
	return a, nil
}
func (s *Store) ListAssets(_ context.Context, t string, l int) ([]domain.Asset, error) {
	o := []domain.Asset{}
	for _, a := range s.state.Assets {
		if t == "" || a.TenantID == t {
			o = append(o, a)
			if l > 0 && len(o) >= l {
				break
			}
		}
	}
	return o, nil
}
func (s *Store) CreatePipeline(_ context.Context, p domain.Pipeline) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Pipelines[p.ID]; ok {
		return domain.ErrConflict
	}
	s.state.Pipelines[p.ID] = p
	return s.persist()
}
func (s *Store) GetPipeline(_ context.Context, id string) (domain.Pipeline, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.state.Pipelines[id]
	if !ok {
		return domain.Pipeline{}, domain.ErrNotFound
	}
	return p, nil
}
func (s *Store) ListPipelines(_ context.Context, t string) ([]domain.Pipeline, error) {
	o := []domain.Pipeline{}
	for _, p := range s.state.Pipelines {
		if t == "" || p.TenantID == t {
			o = append(o, p)
		}
	}
	return o, nil
}
func (s *Store) CreateJob(_ context.Context, j domain.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Jobs[j.ID]; ok {
		return domain.ErrConflict
	}
	s.state.Jobs[j.ID] = j
	return s.persist()
}
func (s *Store) GetJob(_ context.Context, id string) (domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.state.Jobs[id]
	if !ok {
		return domain.Job{}, domain.ErrNotFound
	}
	return j, nil
}
func (s *Store) UpdateJob(_ context.Context, j domain.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Jobs[j.ID]; !ok {
		return domain.ErrNotFound
	}
	s.state.Jobs[j.ID] = j
	return s.persist()
}
func (s *Store) ListJobs(_ context.Context, t string) ([]domain.Job, error) {
	o := []domain.Job{}
	for _, j := range s.state.Jobs {
		if t == "" || j.TenantID == t {
			o = append(o, j)
		}
	}
	return o, nil
}
func (s *Store) CreateArtifact(_ context.Context, a domain.Artifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Artifacts[a.ID]; ok {
		return domain.ErrConflict
	}
	s.state.Artifacts[a.ID] = a
	return s.persist()
}
func (s *Store) ListArtifacts(_ context.Context, j string) ([]domain.Artifact, error) {
	o := []domain.Artifact{}
	for _, a := range s.state.Artifacts {
		if j == "" || a.JobID == j {
			o = append(o, a)
		}
	}
	return o, nil
}
