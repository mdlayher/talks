package main

import (
	"io"
	"strings"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	// START TESTTABLE OMIT
	tests := []struct {
		name string
		r    io.Reader
		out  string
		err  error
	}{
		{
			name: "missing BOOT_IMAGE",
			r:    strings.NewReader(`root=UUID=foo ro`),
			err:  errNoBootImage,
		},
		{
			name: "OK",
			r:    strings.NewReader(`BOOT_IMAGE=/boot/vmlinuz root=UUID=foo ro`),
			out:  `BOOT_IMAGE=/boot/vmlinuz root=UUID=foo ro`,
		},
	}
	// END TESTTABLE OMIT

	// START TESTLOOP OMIT
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Could be strings.NewReader or any other io.Reader
			out, err := ParseCommandLine(tt.r)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %#v\n-  got: %#v",
					want, got)
			}

			if want, got := tt.out, out; want != got {
				t.Fatalf("unexpected command line:\n- want: %q\n-  got: %q",
					want, got)
			}
		})
	}
	// END TESTLOOP OMIT
}
