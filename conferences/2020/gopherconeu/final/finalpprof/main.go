package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/mdlayher/talks/conferences/2020/gopherconeu/final/schedgroup"
)

func main() {
	go http.ListenAndServe(":8081", nil)

	// START OMIT
	ctx, cancel := context.WithCancel(context.Background())

	sg := schedgroup.New(ctx)
	var count uint32
	for i := 0; i < 1_000_000; i++ {
		sg.Delay(10*time.Second, func() {
			atomic.AddUint32(&count, 1)
		})
	}

	if err := sg.Wait(); err != nil {
		log.Fatalf("failed to wait: %v", err)
	}
	cancel()

	log.Printf("done: %d!", atomic.LoadUint32(&count))
	// END OMIT
}
