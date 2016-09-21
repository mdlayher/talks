package main

import "fmt"

func main() {
	a := 4
	b := 2

	c := Divide(a, b)
	fmt.Printf("%d / %d = %d\n", a, b, c)
}

func Divide(a, b int) int {
	return a / b
}
