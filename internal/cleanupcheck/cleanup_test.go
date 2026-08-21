package cleanupcheck_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
	"example.com/media-workflow/internal/storage"
)

func TestCleanupReportsDeletionErrors(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "old")
	if err := os.WriteFile(p, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	_ = os.Chtimes(p, old, old)
	r := storage.Cleanup(d, time.Minute)
	if r.Scanned != 1 || r.Removed != 1 || len(r.Errors) != 0 {
		t.Fatalf("report=%+v", r)
	}
	path := filepath.Join(t.TempDir(), "lock")
	l, err := storage.AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}
	if storage.CleanupLocks([]storage.FileLock{l}) != 1 {
		t.Fatal("lock cleanup count lost")
	}
}

func TestCleanupErrorIsRecorded(t *testing.T){ r:=storage.CleanupReport{}; storage.RecordCleanupError(&r,errors.New("denied")); if len(r.Errors)!=1 {t.Fatal("error lost")} }
func TestRemovePathsCountsCompleted(t *testing.T){ d:=t.TempDir(); p:=filepath.Join(d,"x"); _=os.WriteFile(p,[]byte("x"),0600); if storage.RemovePaths([]string{p,filepath.Join(d,"missing")})!=1 {t.Fatal("wrong removed count")} }
func TestLockReleaseStateIsVisible(t *testing.T){ l,err:=storage.AcquireLock(filepath.Join(t.TempDir(),"lock")); if err!=nil {t.Fatal(err)}; _=l.Release(); if !l.IsReleased(){t.Fatal("release not visible")} }
