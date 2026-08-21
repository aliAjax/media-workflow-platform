package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

type UploadPart struct {
	Number   int
	Offset   int64
	Size     int64
	Checksum string
	Data     []byte
}
type UploadSession struct {
	ID           string
	AssetID      string
	ExpectedSize int64
	Parts        map[int]UploadPart
	Completed    bool
}

func NewUploadSession(id, asset string, size int64) (*UploadSession, error) {
	if id == "" || asset == "" || size < 0 {
		return nil, fmt.Errorf("%w: invalid upload", ErrInvalidInput)
	}
	return &UploadSession{ID: id, AssetID: asset, ExpectedSize: size, Parts: map[int]UploadPart{}}, nil
}
func (u *UploadSession) AddPart(p UploadPart) error {
	if p.Number < 1 || p.Size < 0 || int64(len(p.Data)) != p.Size {
		return fmt.Errorf("%w: invalid part", ErrInvalidInput)
	}
	sum := sha256.Sum256(p.Data)
	checksum := hex.EncodeToString(sum[:])
	if p.Checksum != "" && p.Checksum != checksum {
		return fmt.Errorf("%w: checksum mismatch", ErrConflict)
	}
	p.Checksum = checksum
	if old, ok := u.Parts[p.Number]; ok && old.Checksum != p.Checksum {
		return fmt.Errorf("%w: part conflict", ErrConflict)
	}
	u.Parts[p.Number] = p
	return nil
}
func (u *UploadSession) Assemble() ([]byte, error) {
	keys := make([]int, 0, len(u.Parts))
	for k := range u.Parts {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var out []byte
	var size int64
	for i, k := range keys {
		if k != i+1 {
			return nil, fmt.Errorf("%w: missing part", ErrInvalidInput)
		}
		out = append(out, u.Parts[k].Data...)
		size += u.Parts[k].Size
	}
	if size != u.ExpectedSize {
		return nil, fmt.Errorf("%w: expected %d bytes got %d", ErrInvalidInput, u.ExpectedSize, size)
	}
	u.Completed = true
	return out, nil
}
