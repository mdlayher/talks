//+build !linux

package statfs

import (
	"fmt"
	"runtime"
)

// get is unimplemented.
func get(path string) (*Filesystem, error) {
	return nil, fmt.Errorf("statfs not implemented on %s", runtime.GOOS)
}
