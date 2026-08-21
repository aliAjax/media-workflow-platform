package storage

import (
	"example.com/media-workflow/internal/domain"
	"fmt"
	"os"
	"syscall"
)

type FileLock struct {
	file     *os.File
	released bool
}

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
func (l *FileLock) Release() error {
	if l == nil || l.file == nil {
		return nil
	}
	e := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	if cerr := l.file.Close(); e == nil {
		e = cerr
	}
	l.released = true
	return e
}
func (l *FileLock) IsReleased() bool { return l != nil && l.released }
