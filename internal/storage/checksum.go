package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"example.com/media-workflow/internal/domain"
	"fmt"
	"hash"
	"io"
	"os"
)

func FileChecksum(path string) (string, int64, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", 0, e
	}
	defer f.Close()
	h := sha256.New()
	n, e := io.Copy(h, f)
	if e != nil {
		return "", n, e
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}
func VerifyBytes(expected string, b []byte) error {
	h := sha256.Sum256(b)
	if expected != "" && expected != hex.EncodeToString(h[:]) {
		return fmt.Errorf("%w: checksum", domain.ErrConflict)
	}
	return nil
}
func CopyWithHash(dst io.Writer, src io.Reader) (string, int64, error) {
	h := sha256.New()
	n, e := io.Copy(io.MultiWriter(dst, h), src)
	if e != nil {
		return "", n, e
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}
func NewHash() hash.Hash { return sha256.New() }
