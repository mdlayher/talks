package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	go http.ListenAndServe(":8081", nil)

	// START OMIT
	// Wait for a large number of tasks to complete.
	var wg sync.WaitGroup
	wg.Add(1_000_000)

	ctx, cancel := context.WithCancel(context.Background())
	var count uint32
	for i := 0; i < 1_000_000; i++ {
		// Goroutines either increment the counter or are canceled.
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
			case <-time.After(1 * time.Second):
				atomic.AddUint32(&count, 1)
			}
		}()
	}

	cancel()
	wg.Wait()
	log.Printf("done: %d!", atomic.LoadUint32(&count))
	// END OMIT
}
