package main

import (
	"fmt"
	"time"
)

func main() {
	select {
	case n := <-longOperation():
		fmt.Println("n:", n)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	}
}

func longOperation() chan int {
	c := make(chan int, 0)
	go func() {
		time.Sleep(10 * time.Second)
		c <- 1
	}()
	return c
}
