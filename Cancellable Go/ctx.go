package main

import (
	"context"
	"fmt"
	"time"
)

func doRandom(ctx context.Context) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				out <- 69
			}
		}
	}()
	return out
}

func main() {
	start := time.Now()
	ctx := context.Background()
	Ctx, cancelCtx := context.WithCancel(ctx)
	listen := doRandom(Ctx)
	go func() {
		<-time.After(1 * time.Second)
		cancelCtx()
	}()
	for val := range listen {
		fmt.Println(val)
	}
	fmt.Println(time.Since(start))

}
