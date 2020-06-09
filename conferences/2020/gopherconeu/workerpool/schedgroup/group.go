package schedgroup

import (
	"context"
	"sync"
	"time"
)

// START 1 OMIT

// A Group is a goroutine worker pool which delays tasks for a specified time.
type Group struct {
	taskC chan task
	wg    sync.WaitGroup
}

// New creates a new Group which will use ctx for cancelation.
func New(ctx context.Context) *Group {
	// Spin up n worker goroutines to consume tasks off of taskC.
	const n = 32
	g := &Group{taskC: make(chan task, n)}

	for i := 0; i < n; i++ {
		go g.worker(ctx)
	}

	return g
}

// END 1 OMIT
// START 2 OMIT

// Delay schedules a function to run at or after the specified delay.
func (g *Group) Delay(delay time.Duration, fn func()) {
	g.wg.Add(1)
	g.taskC <- task{
		Delay: delay,
		Call:  fn,
	}
}

// Wait waits for the completion of all scheduled tasks.
func (g *Group) Wait() error {
	g.wg.Wait()
	return nil
}

// A task is a function which is called after the specified delay.
type task struct {
	Delay time.Duration
	Call  func()
}

// END 2 OMIT
// START 3 OMIT

// worker runs a loop which will consume tasks until ctx is canceled.
func (g *Group) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-g.taskC:
			g.work(ctx, t)
		}
	}
}

// work executes a task after a delay or returns if ctx is canceled.
func (g *Group) work(ctx context.Context, t task) {
	defer g.wg.Done()

	select {
	case <-ctx.Done():
	case <-time.After(t.Delay):
		t.Call()
	}
}

// END 3 OMIT
