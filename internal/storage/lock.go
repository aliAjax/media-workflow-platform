package storage

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"os"
	"syscall"
)

type FileLock struct{ file *os.File }

func AcquireLock(path string) (FileLock, error) {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return FileLock{}, e
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		_ = f.Close()
		return FileLock{}, fmt.Errorf("%w: lock unavailable", domain.ErrConflict)
	}
	return FileLock{file: f}, nil
}
func (l FileLock) Release() error {
	if l.file == nil {
		return nil
	}
	_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	return l.file.Close()
}
func (l FileLock) IsReleased() bool { return false }
