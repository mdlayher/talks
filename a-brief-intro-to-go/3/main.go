package main

import "fmt"

func main() {
	// Array. Fixed length.
	a := [3]int{1, 2, 3}

	// Slice! Much more common in Go.
	s := []int{1, 2, 3}

	fmt.Println("array:", a, "slice:", s)

	// Special built-in function append can grow slices.
	s = append(s, 4)
	fmt.Println("slice:", s)
}
