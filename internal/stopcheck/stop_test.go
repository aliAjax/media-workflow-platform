package stopcheck_test

import (
	"testing"
	"example.com/media-workflow/internal/scheduler"
)

func TestStopOnceNoPanic(t *testing.T) {
	r := scheduler.New(nil, nil, nil, nil, nil, nil, 1)
	r.Stop()
	r.Stop()
}
func TestStopStateIsObservable(t *testing.T) { r:=scheduler.New(nil,nil,nil,nil,nil,nil,1); r.Stop(); if !r.Stopped() { t.Fatal("stop state not visible") } }
