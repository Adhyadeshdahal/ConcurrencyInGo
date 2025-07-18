package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// THIS IS ALSO HELPER
func orDone(ctx context.Context, ch <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
				case out <- val:
				}
			}
		}

	}()

	return out
}

func bridge(ctx context.Context, channels <-chan <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var wg sync.WaitGroup
		defer close(out)
		for ch := range channels {
			wg.Add(1)
			nch := orDone(ctx, ch)
			go func(npch <-chan interface{}) {
				defer wg.Done()
				for val := range npch {
					select {
					case <-ctx.Done():
						return
					case out <- val:
					}
				}
			}(nch)
		}
		wg.Wait()
	}()
	return out
}

// THIS IS HELPER
func generator() <-chan <-chan interface{} {
	outer := make(chan (<-chan interface{}))
	go func() {
		defer close(outer)
		for i := 0; i < 5; i++ {
			inner := make(chan interface{})
			go func(n int, ch chan interface{}) {
				defer close(ch)
				for j := 0; j < 3; j++ {
					ch <- fmt.Sprintf("chan %d: val %d", n, j)
					time.Sleep(100 * time.Millisecond)
				}
			}(i, inner)
			outer <- inner
		}
	}()
	return outer
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	incoming := bridge(ctx, generator())

	// Cancel the context after 1 second
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	for val := range incoming {
		fmt.Println(val)
	}
}
