package main

import (
	"container/heap"
	"fmt"
	"time"
)

func main() {
	// START HEAP OMIT

	// Creates a task with a deadline of UNIX timestamp n.
	mkTask := func(n int64) task {
		return task{
			Deadline: time.Unix(n, 0),
			Call:     func() { fmt.Printf("%d!\n", n) },
		}
	}

	// Push tasks onto the heap in any order.
	var tasks tasks
	heap.Push(&tasks, mkTask(500))
	heap.Push(&tasks, mkTask(50))
	heap.Push(&tasks, mkTask(250))
	heap.Push(&tasks, mkTask(0))
	heap.Push(&tasks, mkTask(1000))

	// Tasks will be popped off the heap and called in shortest deadline order.
	for tasks.Len() > 0 {
		t := heap.Pop(&tasks).(task)
		t.Call()
	}

	// END HEAP OMIT
}

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
