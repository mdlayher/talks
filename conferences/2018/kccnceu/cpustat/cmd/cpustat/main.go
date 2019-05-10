// Command cpustat provides basic Linux CPU utilization statistics.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdlayher/talks/conferences/2018/kccnceu/cpustat"
)

func main() {
	f, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatalf("failed to open /proc/stat: %v", err)
	}
	defer f.Close()

	stats, err := cpustat.Scan(f)
	if err != nil {
		log.Fatalf("failed to scan: %v", err)
	}

	for _, s := range stats {
		fmt.Printf("%4s: user: %06d, system: %06d, idle: %06d\n",
			s.ID, s.User, s.System, s.Idle)
	}
}
