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

// defaultBlobRoot returns a stable, writable directory when no explicit root is
// configured. UserCacheDir persists across runs; TempDir is the fallback if the
// cache directory cannot be resolved (e.g. $HOME unset).
func defaultBlobRoot() string {
	if d, err := os.UserCacheDir(); err == nil {
		return filepath.Join(d, "media-workflow", "blobs")
	}
	return filepath.Join(os.TempDir(), "media-workflow-blobs")
}

// NewLocalBlobStore creates a blob store rooted at r. When r is empty the store
// falls back to a usable default directory so first-run asset creation does not
// target an empty path. An explicit r is always honored as-is.
func NewLocalBlobStore(r string) *LocalBlobStore {
	if r == "" {
		r = defaultBlobRoot()
	}
	return &LocalBlobStore{root: r}
}
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
