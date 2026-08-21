package scheduler

import (
	"example.com/media-workflow/internal/domain"
	"sync"
	"testing"
)

func exercisePeek(q *FairQueue, start <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start
	for i := 0; i < 4000; i++ {
		q.Push(domain.Job{ID: "job"})
		q.Peek()
	}
}

func TestConcurrentPeekUsesQueueLock(t *testing.T) {
	q := NewFairQueue()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go exercisePeek(q, start, &wg)
	go exercisePeek(q, start, &wg)
	close(start)
	wg.Wait()
}

func exerciseReset(q *FairQueue, start <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start
	for i := 0; i < 4000; i++ {
		q.Push(domain.Job{ID: "job"})
		q.Reset()
	}
}

func TestConcurrentResetUsesQueueLock(t *testing.T) {
	q := NewFairQueue()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go exerciseReset(q, start, &wg)
	go exerciseReset(q, start, &wg)
	close(start)
	wg.Wait()
}

func TestConcurrentEmptyPopReturnsZero(t *testing.T) {
	q := NewFairQueue()
	if got := q.PopOrZero(); got.ID != "" {
		t.Fatalf("unexpected job %#v", got)
	}
}

func TestConcurrentRetryLimitStopsAtMaximum(t *testing.T) {
	if (RetryPolicy{MaxAttempts: 1}).CanRetry(1) {
		t.Fatal("retry exceeded max attempts")
	}
}
