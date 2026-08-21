package pipeline

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"sort"
)

type CompiledStage struct {
	Stage    domain.Stage
	Index    int
	Children []string
}
type CompiledPipeline struct {
	ID          string
	Version     int
	Stages      []CompiledStage
	TotalWeight int
}

func Compile(p domain.Pipeline) (CompiledPipeline, error) {
	if e := Validate(p); e != nil {
		return CompiledPipeline{}, e
	}
	order, _ := domain.TopologicalOrder(p.Stages)
	children := map[string][]string{}
	for _, s := range p.Stages {
		for _, d := range s.DependsOn {
			children[d] = append(children[d], s.ID)
		}
	}
	out := CompiledPipeline{ID: p.ID, Version: p.Version}
	for i, s := range order {
		w := s.Weight
		if w == 0 {
			w = 1
		}
		out.TotalWeight += w
		sort.Strings(children[s.ID])
		out.Stages = append(out.Stages, CompiledStage{Stage: s, Index: i, Children: children[s.ID]})
	}
	if len(out.Stages) == 0 {
		return out, fmt.Errorf("empty pipeline")
	}
	return out, nil
}
func (c CompiledPipeline) Stage(id string) (CompiledStage, bool) {
	for _, s := range c.Stages {
		if s.Stage.ID == id {
			return s, true
		}
	}
	return CompiledStage{}, false
}
