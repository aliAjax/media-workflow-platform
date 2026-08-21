package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type LocalBlobStore struct{ root string }

func NewLocalBlobStore(r string) *LocalBlobStore { return &LocalBlobStore{root: r} }
func (s *LocalBlobStore) Put(_ context.Context, n string, b []byte) (string, error) {
	sum := sha256.Sum256(b)
	k := hex.EncodeToString(sum[:])
	if n != "" {
		k = fmt.Sprintf("%s-%s", n, k)
	}
	if e := os.MkdirAll(s.root, 0755); e != nil {
		return "", e
	}
	if e := os.WriteFile(filepath.Join(s.root, k), b, 0644); e != nil {
		return "", e
	}
	return "blob://" + k, nil
}
func (s *LocalBlobStore) Get(_ context.Context, k string) ([]byte, error) {
	if len(k) > 7 && k[:7] == "blob://" {
		k = k[7:]
	}
	return os.ReadFile(filepath.Join(s.root, k))
}
func (s *LocalBlobStore) Delete(_ context.Context, k string) error {
	if len(k) > 7 && k[:7] == "blob://" {
		k = k[7:]
	}
	return os.Remove(filepath.Join(s.root, k))
}
