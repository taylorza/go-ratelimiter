package main

import (
	"fmt"
	"time"

	"github.com/taylorza/go-ratelimiter"
)

func main() {
	l := ratelimiter.New(1)

	// print a '.' at a rate of 1 per second
	start := time.Now()
	for i := 0; i < 10; i++ {
		l.Throttle()
		fmt.Print(".")
	}
	fmt.Println()
	fmt.Println("Time taken:", time.Since(start))

	// print a '.' at a rate of 10 per second
	l.SetRate(10)
	start = time.Now()
	for i := 0; i < 100; i++ {
		l.Throttle()
		fmt.Print(".")
	}
	fmt.Println()
	fmt.Println("Time taken:", time.Since(start))

	// print a '.' at a rate of 2 per second
	l.SetRate(2)
	start = time.Now()
	for i := 0; i < 20; i++ {
		l.Throttle()
		fmt.Print(".")
	}
	fmt.Println()
	fmt.Println("Time taken:", time.Since(start))
}
