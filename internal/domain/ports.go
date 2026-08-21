package domain

import "context"

type AssetRepository interface {
	CreateAsset(context.Context, Asset) error
	GetAsset(context.Context, string) (Asset, error)
	ListAssets(context.Context, string, int) ([]Asset, error)
}
type PipelineRepository interface {
	CreatePipeline(context.Context, Pipeline) error
	GetPipeline(context.Context, string) (Pipeline, error)
	ListPipelines(context.Context, string) ([]Pipeline, error)
}
type JobRepository interface {
	CreateJob(context.Context, Job) error
	GetJob(context.Context, string) (Job, error)
	UpdateJob(context.Context, Job) error
	ListJobs(context.Context, string) ([]Job, error)
}
type ArtifactRepository interface {
	CreateArtifact(context.Context, Artifact) error
	ListArtifacts(context.Context, string) ([]Artifact, error)
}
type BlobStore interface {
	Put(context.Context, string, []byte) (string, error)
	Get(context.Context, string) ([]byte, error)
	Delete(context.Context, string) error
}
