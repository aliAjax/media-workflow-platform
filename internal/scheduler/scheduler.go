package scheduler

import (
	"context"
	"example.com/media-workflow/internal/domain"
	"example.com/media-workflow/internal/executor"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Runner struct {
	jobs      domain.JobRepository
	assets    domain.AssetRepository
	pipelines domain.PipelineRepository
	artifacts domain.ArtifactRepository
	blobs     domain.BlobStore
	exec      *executor.Executor
	logger    *slog.Logger
	queue     chan string
	workers   int
	wg        sync.WaitGroup
	stop      chan struct{}
}

func New(j domain.JobRepository, a domain.AssetRepository, p domain.PipelineRepository, ar domain.ArtifactRepository, b domain.BlobStore, l *slog.Logger, w int) *Runner {
	if w < 1 {
		w = 1
	}
	return &Runner{jobs: j, assets: a, pipelines: p, artifacts: ar, blobs: b, exec: executor.New(), logger: l, queue: make(chan string, 256), workers: w, stop: make(chan struct{})}
}
func (r *Runner) Start() {
	for i := 0; i < r.workers; i++ {
		r.wg.Add(1)
		go r.loop(i)
	}
}
func (r *Runner) Stop() { close(r.stop); r.wg.Wait() }
func (r *Runner) Enqueue(id string) error {
	select {
	case r.queue <- id:
		return nil
	default:
		return fmt.Errorf("scheduler queue full")
	}
}
func (r *Runner) loop(_ int) {
	defer r.wg.Done()
	for {
		select {
		case id := <-r.queue:
			r.run(id)
		case <-r.stop:
			return
		}
	}
}
func (r *Runner) run(id string) {
	c := context.Background()
	j, e := r.jobs.GetJob(c, id)
	if e != nil {
		return
	}
	a, e := r.assets.GetAsset(c, j.AssetID)
	if e != nil {
		return
	}
	p, e := r.pipelines.GetPipeline(c, j.PipelineID)
	if e != nil {
		return
	}
	if e = j.Transition(domain.JobRunning); e != nil {
		return
	}
	j.State = domain.JobRunning
	j.Attempt++
	j.UpdatedAt = time.Now()
	_ = r.jobs.UpdateJob(c, j)
	order, e := domain.TopologicalOrder(p.Stages)
	if e != nil {
		return
	}
	total := 0
	weights := map[string]int{}
	for _, s := range order {
		w := s.Weight
		if w == 0 {
			w = 1
		}
		weights[s.ID] = w
		total += w
	}
	done := 0
	for _, s := range order {
		j.CurrentStage = s.ID
		_ = r.jobs.UpdateJob(c, j)
		res, e := r.exec.Run(c, s, a)
		if e != nil {
			j.State = domain.JobFailed
			j.Error = e.Error()
			_ = r.jobs.UpdateJob(c, j)
			return
		}
		uri, e := r.blobs.Put(c, id+"-"+s.ID, res.Data)
		if e == nil {
			_ = r.artifacts.CreateArtifact(c, domain.Artifact{ID: id + "-" + s.ID, JobID: id, AssetID: a.ID, Kind: res.Kind, ContentType: res.ContentType, Checksum: res.Metadata["sha256"], Size: int64(len(res.Data)), URI: uri, CreatedAt: time.Now()})
		}
		done += weights[s.ID]
		j.Progress = done * 100 / total
		_ = r.jobs.UpdateJob(c, j)
	}
	j.State = domain.JobSucceeded
	j.CurrentStage = ""
	j.Progress = 100
	j.UpdatedAt = time.Now()
	_ = r.jobs.UpdateJob(c, j)
	r.logger.Info("job completed", "job_id", id)
}
