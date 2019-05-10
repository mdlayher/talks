package statfs

import (
	"fmt"
)

// A Filesystem contains statistics about a given filesystem.
type Filesystem struct {
	Path  string
	Type  Type
	Files uint64
}

// Get retrieves statistics for the filesystem mounted at path.
func Get(path string) (*Filesystem, error) {
	// Call the OS-specific version of get.
	fs, err := get(path)
	if err != nil {
		return nil, err
	}

	fs.Path = path
	return fs, nil
}

// Type is the type of filesystem detected.
type Type int

// List of possible filesystem types, taken from `man statfs`.
const (
	EXT4 Type = 0xef53
	NFS  Type = 0x6969
	XFS  Type = 0x58465342
)

func (t Type) String() string {
	switch t {
	case EXT4:
		return "EXT4"
	case NFS:
		return "NFS"
	case XFS:
		return "XFS"
	default:
		return fmt.Sprintf("unknown(%x)", int(t))
	}
}
