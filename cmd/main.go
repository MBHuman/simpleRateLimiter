package main

import (
	"context"
	"fmt"
	"ratelimit/internal/ratelimiter"
	"sync"
)

const (
	RATE_LIMIT = 20
)

func main() {
	limiter := ratelimiter.NewLimiterV1(10)

	ch := make(chan interface{}, 100)
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				ch <- "Some data"
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Use a separate goroutine to run limiter.Run, as it blocks until completion
	limiter.Run(RATE_LIMIT, ch, func(ctx context.Context, data interface{}) (res interface{}, err error) {
		fmt.Println(data)
		return
	})

}
