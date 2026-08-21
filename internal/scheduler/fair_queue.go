package scheduler

import (
	"container/heap"
	"example.com/media-workflow/internal/domain"
	"sync"
)

type item struct {
	job   domain.Job
	index int
}
type priorityQueue []*item

func (q priorityQueue) Len() int { return len(q) }
func (q priorityQueue) Less(i, j int) bool {
	if q[i].job.Progress == q[j].job.Progress {
		return q[i].job.CreatedAt.Before(q[j].job.CreatedAt)
	}
	return q[i].job.Progress < q[j].job.Progress
}
func (q priorityQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].index = i; q[j].index = j }
func (q *priorityQueue) Push(v any)   { x := v.(*item); x.index = len(*q); *q = append(*q, x) }
func (q *priorityQueue) Pop() any     { old := *q; n := len(old); x := old[n-1]; *q = old[:n-1]; return x }

type FairQueue struct {
	mu sync.Mutex
	q  priorityQueue
}

func NewFairQueue() *FairQueue { q := priorityQueue{}; heap.Init(&q); return &FairQueue{q: q} }
func (f *FairQueue) Push(j domain.Job) {
	f.mu.Lock()
	defer f.mu.Unlock()
	heap.Push(&f.q, &item{job: j})
}
func (f *FairQueue) Pop() (domain.Job, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.q) == 0 {
		return domain.Job{}, false
	}
	return heap.Pop(&f.q).(*item).job, true
}
func (f *FairQueue) Len() int { f.mu.Lock(); defer f.mu.Unlock(); return len(f.q) }
