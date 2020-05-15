package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// START OMIT
	// Create a WaitGroup to wait for goroutine completion.
	const num = 5
	var wg sync.WaitGroup
	wg.Add(num)

	// Schedule tasks to run every 500ms which print a number to the screen.
	for i := 0; i < num; i++ {
		// 0ms, 500ms, 1000ms, etc.
		n := i + 1
		time.AfterFunc(500*time.Millisecond, func() {
			defer wg.Done()
			log.Println(n)
		})
	}

	// Wait for all of the scheduled tasks to complete.
	wg.Wait()
	log.Println("done!")
	// END OMIT
}
