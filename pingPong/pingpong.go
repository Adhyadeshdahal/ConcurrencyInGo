package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	done := make(chan interface{})

	go func() {
		ch1, ch2 := make(chan interface{}), make(chan interface{})
		defer close(ch1)
		defer close(ch2)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch1:
					fmt.Println("ping")
					ch2 <- 0
				case <-done:
					return
				}
			}

		}()
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch2:
					fmt.Println("pong")
					ch1 <- 1
				case <-done:
					return
				}
			}

		}()
		ch1 <- "start"
		wg.Wait()
		// return done
	}()
	time.Sleep(5 * time.Second)
	done <- "finished"
	return
}
