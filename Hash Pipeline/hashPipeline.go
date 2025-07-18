package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func generator(ctx context.Context, dir string) <-chan string {
	out := make(chan string, 128)
	go func() {
		defer close(out)
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("readdir %s: %v\n", dir, err)
			return
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- filepath.Join(dir, e.Name()):
			}
		}
	}()
	return out
}

func worker(ctx context.Context, paths <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case p, ok := <-paths:
			if !ok {
				return
			}
			f, err := os.ReadFile(p)
			if err != nil {
				log.Printf("open %s: %v\n", p, err)
				continue
			}
			h := sha256.New()
			h.Write([]byte(f))
			results <- fmt.Sprintf("hash of %s is %x", p, h.Sum(nil))
		}
	}
}

func calculateHash(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string, 128)
	workers := runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, in, out, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ctx := context.Background()

	start := time.Now()
	for res := range calculateHash(ctx, generator(ctx, `D:\downloads\`)) {
		fmt.Println(res)
	}
	fmt.Println("It took", time.Since(start))
}
