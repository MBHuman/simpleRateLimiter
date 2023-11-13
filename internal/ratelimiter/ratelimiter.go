package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type ProcessData func(ctx context.Context, data interface{}) (interface{}, error)

type ILimiter interface {
	Run(limit int, data chan interface{}, processData ProcessData)
}

type LimiterV1 struct{}

func NewLimiterV1() ILimiter {
	return &LimiterV1{}
}

func (l *LimiterV1) Run(limit int, data chan interface{}, processData ProcessData) {
	var wg sync.WaitGroup
	ctx := context.TODO()
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tick := time.NewTicker(time.Second / time.Duration(limit))

			for e := range data {
				processData(ctx, e)
				<-tick.C
			}
		}()
	}
	wg.Wait()
}
