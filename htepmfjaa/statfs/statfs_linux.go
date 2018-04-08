//+build linux

package statfs

import "golang.org/x/sys/unix"

// get is the Linux implementation of get.
func get(path string) (*Filesystem, error) {
	// Structure is populated by kernel by passing a pointer to it.
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return nil, err
	}

	// Return the non-OS-specific structure to support multiple OSes.
	return &Filesystem{
		Type:  Type(s.Type),
		Files: s.Files,
	}, nil
}
