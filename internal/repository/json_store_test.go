package repository

import (
	"context"
	"example.com/media-workflow/internal/domain"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentJobListAndUpdate(t *testing.T) {
	s, err := New(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateJob(context.Background(), domain.Job{ID: "job-1", TenantID: "tenant", State: domain.JobQueued}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 300; i++ {
			_, _ = s.ListJobs(context.Background(), "tenant")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 300; i++ {
			_ = s.UpdateJob(context.Background(), domain.Job{ID: "job-1", TenantID: "tenant", State: domain.JobRunning, Progress: i % 101})
		}
	}()
	close(start)
	wg.Wait()
}

func TestConcurrentAssetListAndCreate(t *testing.T) {
	s, err := New(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_, _ = s.ListAssets(context.Background(), "tenant", 0)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_ = s.CreateAsset(context.Background(), domain.Asset{ID: "a-" + string(rune(i)), TenantID: "tenant", Name: "n"})
		}
	}()
	close(start)
	wg.Wait()
}

func TestConcurrentPipelineListAndCreate(t *testing.T) {
	s, err := New(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_, _ = s.ListPipelines(context.Background(), "tenant")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_ = s.CreatePipeline(context.Background(), domain.Pipeline{ID: "p-" + string(rune(i)), TenantID: "tenant", Name: "n"})
		}
	}()
	close(start)
	wg.Wait()
}

func TestConcurrentArtifactListAndCreate(t *testing.T) {
	s, err := New(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_, _ = s.ListArtifacts(context.Background(), "job")
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 200; i++ {
			_ = s.CreateArtifact(context.Background(), domain.Artifact{ID: "x-" + string(rune(i)), JobID: "job"})
		}
	}()
	close(start)
	wg.Wait()
}
