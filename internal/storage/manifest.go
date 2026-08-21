package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"example.com/media-workflow/internal/domain"
	"fmt"
	"sort"
)

type Manifest struct {
	JobID     string            `json:"job_id"`
	Artifacts []domain.Artifact `json:"artifacts"`
	Digest    string            `json:"digest"`
}

func BuildManifest(job string, items []domain.Artifact) Manifest {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	m := Manifest{JobID: job, Artifacts: items}
	b, _ := json.Marshal(m)
	sum := sha256.Sum256(b)
	m.Digest = hex.EncodeToString(sum[:])
	return m
}
func VerifyManifest(m Manifest) error {
	if m.JobID == "" || m.Digest == "" {
		return fmt.Errorf("manifest incomplete")
	}
	d := m.Digest
	m.Digest = ""
	b, _ := json.Marshal(m)
	sum := sha256.Sum256(b)
	if hex.EncodeToString(sum[:]) != d {
		return fmt.Errorf("manifest digest mismatch")
	}
	return nil
}

func SelectArtifacts(items []domain.Artifact, kind string) []domain.Artifact {
	out := items[:0]
	for _, item := range items {
		if kind == "" || item.Kind == kind {
			out = append(out, item)
		}
	}
	return out
}
