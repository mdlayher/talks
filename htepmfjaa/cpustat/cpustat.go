// Package cpustat provides an example parser for Linux CPU utilization statistics.
package cpustat

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// A CPUStat contains statistics for an individual CPU.
type CPUStat struct {
	// The ID of the CPU.
	ID string

	// The time, in USER_HZ (typically 1/100th of a second),
	// spent in each of user, system, and idle modes.
	User, System, Idle int
}

// Scan reads and parses CPUStat information from r.
func Scan(r io.Reader) ([]CPUStat, error) {
	// Skip the first summarized line.
	s := bufio.NewScanner(r)
	s.Scan()

	var stats []CPUStat
	for s.Scan() {
		// Each CPU stats line should have exactly 11 fields.
		const nFields = 11
		fields := strings.Fields(string(s.Bytes()))
		if len(fields) != nFields {
			continue
		}

		// The values we care about (user, system, idle) lie at indices
		// 1, 3, and 4, respectively.  Parse these into the array.
		var times [3]int
		for i, idx := range []int{1, 3, 4} {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, err
			}

			times[i] = v
		}

		stats = append(stats, CPUStat{
			// First field is the CPU's ID.
			ID:     fields[0],
			User:   times[0],
			System: times[1],
			Idle:   times[2],
		})
	}

	// Be sure to check the error!
	if err := s.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
