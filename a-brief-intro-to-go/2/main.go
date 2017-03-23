package main

import "fmt"

func main() {
	fmt.Println(sum2(1, 2))
	fmt.Println(sumN(1, 2, 3))
}

func sum2(x, y int) int {
	return sumN(x, y)
}

// Varidiac functions possible too!
func sumN(numbers ...int) int {
	var sum int
	for _, n := range numbers {
		sum += n
	}

	return sum
}
