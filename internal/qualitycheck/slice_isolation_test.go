package qualitycheck_test

import (
	"example.com/media-workflow/internal/domain"
	"example.com/media-workflow/internal/storage"
	"testing"
)

func TestSelectArtifactsDoesNotReuseInput(t *testing.T) {
	items := []domain.Artifact{{ID: "a", Kind: "x"}, {ID: "b", Kind: "y"}}
	out := storage.SelectArtifacts(items, "x")
	out[0].ID = "changed"
	if items[0].ID != "a" {
		t.Fatal("selection aliases input")
	}
}
func TestQualityEvaluationCopiesWarnings(t *testing.T) {
	warnings := make([]string, 1, 2)
	warnings[0] = "existing"
	r := domain.QualityPolicy{}.Evaluate(domain.QualityReport{Warnings: warnings})
	r.Warnings[0] = "changed"
	if warnings[0] != "existing" {
		t.Fatal("quality warnings alias input")
	}
}
func TestNormalizeWarningsReturnsCopy(t *testing.T) {
	in := []string{"first"}
	out := domain.NormalizeWarnings(in)
	out[0] = "changed"
	if in[0] != "first" {
		t.Fatal("normalized warnings alias input")
	}
}
