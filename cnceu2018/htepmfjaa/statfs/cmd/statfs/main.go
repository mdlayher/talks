// Command statfs retrieves statistics about filesystems.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mdlayher/talks/htepmfjaa/statfs"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		fmt.Println("usage: statfs [path]")
		return
	}

	fs, err := statfs.Get(path)
	if err != nil {
		log.Fatalf("failed to get filesystem: %v", err)
	}

	fmt.Printf("%s (%s): %d files\n", fs.Path, fs.Type, fs.Files)
}
