package domain

import "time"

type EventType string

const (
	EventAssetCreated    EventType = "asset.created"
	EventPipelineCreated EventType = "pipeline.created"
	EventJobQueued       EventType = "job.queued"
	EventJobProgress     EventType = "job.progress"
	EventJobCompleted    EventType = "job.completed"
	EventJobFailed       EventType = "job.failed"
)

type DomainEvent struct {
	ID          string
	Type        EventType
	AggregateID string
	TenantID    string
	Payload     any
	CreatedAt   time.Time
	Attempt     int
}

func NewEvent(id string, t EventType, aggregate, tenant string, payload any) DomainEvent {
	return DomainEvent{ID: id, Type: t, AggregateID: aggregate, TenantID: tenant, Payload: payload, CreatedAt: time.Now()}
}
func (e DomainEvent) Retry() DomainEvent { e.Attempt++; return e }
