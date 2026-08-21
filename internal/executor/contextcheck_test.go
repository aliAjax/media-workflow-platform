package executor_test

import (
 "context"
 "errors"
 "testing"
 "example.com/media-workflow/internal/domain"
 "example.com/media-workflow/internal/executor"
 "example.com/media-workflow/internal/observability"
)
func TestExecutorSequenceHonorsCancellation(t *testing.T){ c,x:=context.WithCancel(context.Background()); x(); if err:=executor.New().RunSequence(c,[]domain.Stage{{Kind:"probe"}},domain.Asset{ID:"a"}); !errors.Is(err,context.Canceled){t.Fatalf("%v",err)} }
func TestCheckReadyPassesContext(t *testing.T){ c,x:=context.WithCancel(context.Background()); x(); err:=observability.CheckReady(c,observability.Check{Name:"c",Run:func(ctx context.Context)error{return ctx.Err()}}); if !errors.Is(err,context.Canceled){t.Fatalf("%v",err)} }
func TestRunChecksStopsOnCanceledContext(t *testing.T){ c,x:=context.WithCancel(context.Background()); x(); r:=observability.RunChecks(c,[]observability.Check{{Name:"c",Run:func(ctx context.Context)error{return ctx.Err()}}}); if r.Status!="degraded"{t.Fatalf("%+v",r)} }
