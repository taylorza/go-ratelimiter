package main

import (
	"fmt"
	"time"

	"github.com/taylorza/go-ratelimiter"
)

func main() {
	// create a rate limiter that will limit work to 1 per second
	l := ratelimiter.New(1)

	// print a '.' at a rate of 1 per second
	start := time.Now()
	for i := 0; i < 10; i++ {
		l.Throttle()
		fmt.Print(".")
	}
	fmt.Println()
	fmt.Println("Time taken:", time.Since(start))
}
