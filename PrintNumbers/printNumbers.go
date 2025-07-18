package main

import (
	"fmt"
	"time"
)

type ans struct {
	i, id int
}

func Print() {
	for i := 0; i < 2500; i++ {
		fmt.Println(ans{i, 0})
	}
}

// func publisher(ch chan<- ans, id int) {
// 	for i := 0; i < 5; i++ {
// 		ch <- ans{i, id}
// 	}
// }

// func Print() {
// 	channel := func() <-chan ans {
// 		ch := make(chan ans, 500)
// 		go func() {
// 			defer close(ch)
// 			var wg sync.WaitGroup
// 			wg.Add(500)
// 			for i := 0; i < 500; i++ {
// 				go func(id int) {
// 					defer wg.Done()
// 					publisher(ch, id)
// 				}(i)
// 			}
// 			wg.Wait()
// 		}()

// 		return ch
// 	}()
// 	var wg sync.WaitGroup
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			defer wg.Done()
// 			for val := range channel {
// 				fmt.Println(val)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }

func main() {
	start := time.Now()
	Print()
	fmt.Println(time.Since(start))
}
