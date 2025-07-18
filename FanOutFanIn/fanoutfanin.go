package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int) {
	for i := 0; i < 1000; i++ {
		ch <- i
	}
}

func fanOut(ch <-chan int, n int) []chan int {
	output := make([]chan int, n)
	for i := 0; i < n; i++ {
		output[i] = make(chan int)
		go func(index int) {
			defer close(output[index])
			for val := range ch {
				output[index] <- val
			}
		}(i)
	}
	return output
}

func consumer(ch <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for val := range ch {
			output <- (val * val)
		}
	}()
	return output
}

func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		var wg sync.WaitGroup

		for _, ch := range channels {
			wg.Add(1)
			go func(c <-chan int) {
				defer wg.Done()
				for val := range c {
					out <- val
				}
			}(ch)
		}
		wg.Wait()
	}()
	return out
}
func main() {

	ch := make(chan int)
	go func() {
		producer(ch)
		close(ch)
	}()

	fanOutChannels := fanOut(ch, 5)
	consumedChannels := make([]<-chan int, len(fanOutChannels))
	for i, ch := range fanOutChannels {
		consumedChannels[i] = consumer(ch)
	}

	fanInChannel := fanIn(consumedChannels...)

	for val := range fanInChannel {
		fmt.Println(val)
	}
}
