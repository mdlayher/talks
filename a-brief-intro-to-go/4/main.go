package main

import (
	"fmt"
	"math"
)

// Any type with an Area method satisfies this interface.
type Arear interface {
	Area() float64
}

func main() {
	c := Circle{Radius: 2}
	printArea(c)
}

func printArea(a Arear) {
	fmt.Println("area:", a.Area())
}

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return math.Pi * math.Pow(c.Radius, 2) }
