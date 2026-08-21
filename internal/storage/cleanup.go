package storage

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CleanupReport struct {
	Scanned int
	Removed int
	Bytes   int64
	Errors  []string
}

func Cleanup(root string, olderThan time.Duration) CleanupReport {
	r := CleanupReport{}
	cutoff := time.Now().Add(-olderThan)
	_ = filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			r.Errors = append(r.Errors, e.Error())
			return nil
		}
		if info == nil || info.IsDir() || strings.HasSuffix(path, ".keep") {
			return nil
		}
		r.Scanned++
		if info.ModTime().Before(cutoff) {
			if e := os.Remove(path); e != nil {
				r.Errors = append(r.Errors, e.Error())
			} else {
				r.Removed++
				r.Bytes += info.Size()
			}
		}
		return nil
	})
	return r
}

func CleanupLocks(locks []FileLock) int {
	released := 0
	for i := range locks {
		if locks[i].Release() == nil {
			released++
		}
	}
	return released
}

func RecordCleanupError(r *CleanupReport, err error) {
	r.Errors = append(r.Errors, err.Error())
}

func RemovePaths(paths []string) int {
	removed := 0
	for _, path := range paths {
		if e := os.Remove(path); e == nil {
			removed++
		}
	}
	return removed
}
