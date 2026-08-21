package qualitycheck_test

import (
	"example.com/media-workflow/internal/domain"
	"example.com/media-workflow/internal/storage"
	"testing"
)

func TestBuildManifestDoesNotMutateArtifacts(t *testing.T) {
	items := make([]domain.Artifact, 2, 4)
	items[0] = domain.Artifact{ID: "b"}
	items[1] = domain.Artifact{ID: "a"}
	storage.BuildManifest("job", items)
	if items[0].ID != "b" || items[1].ID != "a" {
		t.Fatalf("caller slice reordered: %#v", items)
	}
}
