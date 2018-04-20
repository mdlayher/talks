package main

import (
	"errors"
	"fmt"
)

func main() {
	a := 4
	b := 0

	c, err := Divide(a, b)
	if err != nil {
		fmt.Println("fatal error:", err)
		return
	}

	fmt.Printf("%d / %d = %d\n", a, b, c)
}

// START OMIT
var errDivideByZero = errors.New("division by zero")

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errDivideByZero
	}

	return a / b, nil
}

// END OMIT
