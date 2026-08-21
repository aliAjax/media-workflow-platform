package pagingcheck_test

import (
	"example.com/media-workflow/internal/application"
	"example.com/media-workflow/internal/domain"
	"testing"
)

func TestPaginateJobsReturnsIndependentSlice(t *testing.T) {
	items := make([]domain.Job, 2, 2)
	items[0] = domain.Job{ID: "one"}
	items[1] = domain.Job{ID: "keep"}
	p := application.PaginateJobs(items, "", 1)
	p.Items = append(p.Items, domain.Job{ID: "two"})
	if items[0].ID != "one" || items[1].ID != "keep" {
		t.Fatalf("input mutated: %#v", items)
	}
}

func TestPageCloneKeepsOriginalItems(t *testing.T) {
	p := application.Page[domain.Job]{Items: []domain.Job{{ID: "original"}}}
	cloned := p.Clone()
	cloned.Items[0].ID = "changed"
	if p.Items[0].ID != "original" {
		t.Fatalf("page clone mutated original: %#v", p.Items)
	}
}

func TestCommandBusListIsolated(t *testing.T) {
	b := application.NewCommandBus()
	_ = b.Dispatch(nil, "1", "x", "a")
	out := b.List()
	out[0].ID = "changed"
	if b.List()[0].ID != "1" {
		t.Fatal("command log leaked")
	}
}

func TestCommandBusRecentWindowIsolated(t *testing.T) {
	b := application.NewCommandBus()
	_ = b.Dispatch(nil, "1", "start", "asset")
	_ = b.Dispatch(nil, "2", "finish", "asset")
	out := b.Recent(1)
	out[0].Type = "changed"
	if b.Recent(1)[0].Type != "finish" {
		t.Fatal("recent command window leaked")
	}
}
