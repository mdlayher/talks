package main

import "testing"

func TestDivide(t *testing.T) {
	a := 4
	b := 2

	c := Divide(a, b)
	if c != 2 {
		t.Fatal("c should be 2")
	}
}
