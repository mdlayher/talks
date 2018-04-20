package main

import "fmt"

func main() {
	in := make(chan int, 0)
	go func() {
		for i := 0; i < 10; i++ {
			in <- i
		}
		fmt.Println("send done")
		close(in)
	}()
	fmt.Println("sum:", sum(in))
}

func sum(in <-chan int) int {
	var sum int
	for n := range in {
		sum += n
	}
	fmt.Println("receive done")
	return sum
}
