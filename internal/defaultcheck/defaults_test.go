package defaultcheck_test

import "testing"
import "example.com/media-workflow/internal/storage"
import "example.com/media-workflow/internal/domain"
import "context"

func TestDefaultBlobStoreRootIsUsable(t *testing.T) {
	s := storage.NewLocalBlobStore("")
	if s == nil {
		t.Fatal("nil blob store")
	}
	if _, err := s.Put(context.Background(), "probe", []byte("x")); err != nil {
		t.Fatalf("default root unusable: %v", err)
	}
}

func TestUploadSessionInitializesParts(t *testing.T) { u,err:=domain.NewUploadSession("u","a",1); if err!=nil {t.Fatal(err)}; if err=u.AddPart(domain.UploadPart{Number:1,Size:1,Data:[]byte("x")}); err!=nil {t.Fatal(err)} }
func TestQuotaLedgerInitializesItems(t *testing.T) { q:=domain.NewQuotaLedger(); if err:=q.Set(domain.Quota{TenantID:"t",MaxBytes:1}); err!=nil {t.Fatal(err)} }
func TestResourceClaimLabelsAreCopied(t *testing.T) { c:=domain.ResourceClaim{Labels:map[string]string{"zone":"a"}}; out:=c.CopyLabels(); out["zone"]="b"; if c.Labels["zone"]!="a" {t.Fatal("labels escaped") } }
