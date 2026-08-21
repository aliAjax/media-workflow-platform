package application

import (
	"context"
	"example.com/media-workflow/internal/domain"
	"example.com/media-workflow/internal/pipeline"
	"fmt"
	"time"
)

type Service struct {
	assets    domain.AssetRepository
	pipelines domain.PipelineRepository
	jobs      domain.JobRepository
	artifacts domain.ArtifactRepository
	enqueue   func(string) error
	now       func() time.Time
}

func New(a domain.AssetRepository, p domain.PipelineRepository, j domain.JobRepository, ar domain.ArtifactRepository, e func(string) error) *Service {
	return &Service{assets: a, pipelines: p, jobs: j, artifacts: ar, enqueue: e, now: time.Now}
}
func (s *Service) CreateAsset(c context.Context, a domain.Asset) (domain.Asset, error) {
	if a.ID == "" {
		a.ID = fmt.Sprintf("asset-%d", s.now().UnixNano())
	}
	if a.State == "" {
		a.State = domain.AssetReady
	}
	a.CreatedAt = s.now()
	a.UpdatedAt = a.CreatedAt
	if e := a.Validate(); e != nil {
		return domain.Asset{}, e
	}
	return a, s.assets.CreateAsset(c, a)
}
func (s *Service) CreatePipeline(c context.Context, p domain.Pipeline) (domain.Pipeline, error) {
	if p.ID == "" {
		p.ID = fmt.Sprintf("pipe-%d", s.now().UnixNano())
	}
	if e := pipeline.Validate(p); e != nil {
		return domain.Pipeline{}, e
	}
	if p.Version == 0 {
		p.Version = 1
	}
	p.CreatedAt = s.now()
	return p, s.pipelines.CreatePipeline(c, p)
}
func (s *Service) CreateJob(c context.Context, j domain.Job) (domain.Job, error) {
	if j.ID == "" {
		j.ID = fmt.Sprintf("job-%d", s.now().UnixNano())
	}
	if j.State == "" {
		j.State = domain.JobQueued
	}
	j.CreatedAt = s.now()
	j.UpdatedAt = j.CreatedAt
	if _, e := s.assets.GetAsset(c, j.AssetID); e != nil {
		return domain.Job{}, e
	}
	if _, e := s.pipelines.GetPipeline(c, j.PipelineID); e != nil {
		return domain.Job{}, e
	}
	if e := s.jobs.CreateJob(c, j); e != nil {
		return domain.Job{}, e
	}
	return j, s.enqueue(j.ID)
}
func (s *Service) GetJob(c context.Context, id string) (domain.Job, []domain.Artifact, error) {
	j, e := s.jobs.GetJob(c, id)
	if e != nil {
		return j, nil, e
	}
	a, e := s.artifacts.ListArtifacts(c, id)
	return j, a, e
}
