package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Nothing to see here
}

// START CMDLINE OMIT
func ReadCommandLine() (string, error) {
	f, _ := os.Open("/proc/cmdline")
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	if !strings.Contains(string(b), "BOOT_IMAGE") {
		return "", errors.New("no BOOT_IMAGE parameter found")
	}

	return string(b), nil
}

// END CMDLINE OMIT

// START IOREADER OMIT
func OpenCommandLine() (string, error) {
	// No need to test that os.Open does the right thing.  We just
	// care that our ParseCommandLine() function is correct.
	f, _ := os.Open("/proc/cmdline")
	defer f.Close()
	return ParseCommandLine(f)
}

var errNoBootImage = errors.New("no BOOT_IMAGE parameter found")

func ParseCommandLine(r io.Reader) (string, error) {
	b, _ := ioutil.ReadAll(r)
	if !strings.Contains(string(b), "BOOT_IMAGE") {
		return "", errNoBootImage
	}
	return string(b), nil
}

// END IOREADER OMIT
