package main

import (
	"context"
	"log"
	"time"

	"github.com/mdlayher/schedgroup"
)

func main() {
	// START OMIT
	// Initialize the schedgroup.
	sg := schedgroup.New(context.Background())

	// Schedule tasks to run every 500ms which print a number to the screen.
	for i := 0; i < 5; i++ {
		// 0ms, 500ms, 1000ms, etc.
		n := i + 1
		sg.Delay(time.Duration(n)*500*time.Millisecond, func() {
			log.Println(n)
		})
	}

	// Wait for all of the scheduled tasks to complete.
	if err := sg.Wait(); err != nil {
		log.Fatalf("failed to wait: %v", err)
	}
	log.Println("done!")
	// END OMIT
}
