// Command cpustat_exporter exports Prometheus metrics for basic Linux CPU
// utilization statistics.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mdlayher/talks/conferences/2018/kccnceu/cpustat"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Called on each collector.Collect.
	stats := func() ([]cpustat.CPUStat, error) {
		f, err := os.Open("/proc/stat")
		if err != nil {
			return nil, fmt.Errorf("failed to open /proc/stat: %v", err)
		}
		defer f.Close()

		return cpustat.Scan(f)
	}

	// Make Prometheus client aware of our collector.
	c := newCollector(stats)
	prometheus.MustRegister(c)

	// Set up HTTP handler for metrics.
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Start listening for HTTP connections.
	const addr = ":8888"
	log.Printf("starting cpustat exporter on %q", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("cannot start cpustat exporter: %s", err)
	}
}

var _ prometheus.Collector = &collector{}

// A collector is a prometheus.Collector for Linux CPU stats.
type collector struct {
	// Possible metric descriptions.
	TimeUserHertzTotal *prometheus.Desc

	// A parameterized function used to gather metrics.
	stats func() ([]cpustat.CPUStat, error)
}

// newCollector constructs a collector using a stats function.
func newCollector(stats func() ([]cpustat.CPUStat, error)) prometheus.Collector {
	return &collector{
		TimeUserHertzTotal: prometheus.NewDesc(
			// Name of the metric.
			"cpustat_time_user_hertz_total",
			// The metric's help text.
			"Time in USER_HZ a given CPU spent in a given mode.",
			// The metric's variable label dimensions.
			[]string{"cpu", "mode"},
			// The metric's constant label dimensions.
			nil,
		),

		stats: stats,
	}
}

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	// Gather metadata about each metric.
	ds := []*prometheus.Desc{
		c.TimeUserHertzTotal,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements prometheus.Collector.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	// Take a stats snapshot.  Must be concurrency safe.
	stats, err := c.stats()
	if err != nil {
		// If an error occurs, send an invalid metric to notify
		// Prometheus of the problem.
		ch <- prometheus.NewInvalidMetric(c.TimeUserHertzTotal, err)
		return
	}

	for _, s := range stats {
		tuples := []struct {
			mode string
			v    int
		}{
			{mode: "user", v: s.User},
			{mode: "system", v: s.System},
			{mode: "idle", v: s.Idle},
		}

		for _, t := range tuples {
			// prometheus.Collector implementations should always use
			// "const metric" constructors.
			ch <- prometheus.MustNewConstMetric(
				c.TimeUserHertzTotal,
				prometheus.CounterValue,
				float64(t.v),
				s.ID, t.mode,
			)
		}
	}
}
