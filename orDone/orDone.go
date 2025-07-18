package main

import (
	"fmt"
	"math/rand"
	"time"
)

func orDone(done <-chan interface{}, ch <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case val, ok := <-ch:
				if !ok {
					return
				}

				select {
				case <-done:
				case out <- val:
				}
			}
		}

	}()

	return out
}

func GeneratorToGenerateEven(done <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		i := 0
		defer close(out)
		for {
			select {
			case out <- i:
				i += 2
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
	done := make(chan interface{})
	go func() {
		for {
			if rand.Float64()*100 > 90 {
				done <- 0
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	even := GeneratorToGenerateEven(done)
	newEven := orDone(done, even)
	for nE := range newEven {
		fmt.Println(nE)
	}

}
