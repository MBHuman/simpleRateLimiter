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

type LimiterV1 struct {
	gCount int
}

func NewLimiterV1(gCount int) ILimiter {
	return &LimiterV1{gCount: gCount}
}

func (l *LimiterV1) Run(limit int, data chan interface{}, processData ProcessData) {
	var wg sync.WaitGroup
	ctx := context.TODO()
	tick := time.NewTicker(time.Second / time.Duration(limit))

	for i := 0; i < l.gCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range data {
				<-tick.C
				processData(ctx, e)
			}
		}()
	}
	wg.Wait()
}
