package schedgroup

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

// START GROUP OMIT

// A Group is a goroutine worker pool which schedules tasks to be performed
// after a specified time. A Group must be created with the New constructor.
type Group struct {
	ctx    context.Context
	cancel func()
	mu     sync.Mutex
	tasks  tasks
}

// New creates a new Group which will use ctx for cancelation.
func New(ctx context.Context) *Group {
	// Monitor goroutine context and cancelation.
	mctx, cancel := context.WithCancel(ctx)

	g := &Group{
		ctx:    ctx,
		cancel: cancel,
	}
	go g.monitor(mctx)

	return g
}

// END GROUP OMIT

// START DELAY OMIT

// Delay schedules a function to run at or after the specified delay. Delay
// is a convenience wrapper for Schedule which adds delay to the current time.
func (g *Group) Delay(delay time.Duration, fn func()) {
	g.Schedule(time.Now().Add(delay), fn)
}

// Schedule schedules a function to run at or after the specified time.
func (g *Group) Schedule(when time.Time, fn func()) {
	g.mu.Lock()
	defer g.mu.Unlock()

	heap.Push(&g.tasks, task{
		Deadline: when,
		Call:     fn,
	})
}

// END DELAY OMIT

// START WAIT OMIT

// Wait waits for the completion of all scheduled tasks, or for context cancelation.
func (g *Group) Wait() error {
	t := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-g.ctx.Done():
			// Caller asked for cancelation.
			return g.ctx.Err()
		case <-t.C:
		}
		g.mu.Lock()
		if len(g.tasks) == 0 {
			// No more tasks left, cancel the monitor goroutine.
			defer g.mu.Unlock()
			g.cancel()
			t.Stop()
			return nil
		}
		g.mu.Unlock()
	}
}

// END WAIT OMIT

// START MONITOR OMIT

// monitor triggers tasks at the interval specified by g.Interval until ctx
// is canceled.
func (g *Group) monitor(ctx context.Context) error {
	t := time.NewTicker(1 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			// monitor's cancelation is expected and should not result in an
			// error being returned to the caller.
			return nil
		case now := <-t.C:
			g.trigger(now)
		}
	}
}

// END MONITOR OMIT

// START TRIGGER OMIT

// trigger checks for scheduled tasks and runs them if they are scheduled
// on or after the time specified by now.
func (g *Group) trigger(now time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for g.tasks.Len() > 0 {
		// Check the first task's readiness, but don't pop it off the heap
		// just in case it isn't ready.
		next := &g.tasks[0]
		if next.Deadline.After(now) {
			// Earliest scheduled task is not ready.
			return
		}

		// This task is ready, pop it from the heap and run it.
		t := heap.Pop(&g.tasks).(task)
		go t.Call()
	}
}

// END TRIGGER OMIT

// START TASK OMIT

// A task is a function which is called after the specified deadline.
type task struct {
	Deadline time.Time
	Call     func()
}

// tasks implements heap.Interface.
type tasks []task

var _ heap.Interface = &tasks{}

func (pq tasks) Len() int { return len(pq) }

func (pq tasks) Less(i, j int) bool { return pq[i].Deadline.Before(pq[j].Deadline) }

func (pq tasks) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *tasks) Push(x interface{}) { *pq = append(*pq, x.(task)) }

func (pq *tasks) Pop() (item interface{}) {
	n := len(*pq)
	item, *pq = (*pq)[n-1], (*pq)[:n-1]
	return item
}

// END TASK OMIT
