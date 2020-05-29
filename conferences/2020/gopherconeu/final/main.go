package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/mdlayher/talks/conferences/2020/gopherconeu/final/schedgroup"
)

func main() {
	go http.ListenAndServe(":8081", nil)

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
