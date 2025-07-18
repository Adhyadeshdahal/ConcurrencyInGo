package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func Producer(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case out <- rand.Float64():
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	ch := Producer(ctx)
	go func() {
		<-time.After(2 * time.Second)
		cancel()
	}()
	for e := range ch {
		fmt.Println(e)
	}

}
