package main

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/nfnt/resize"
)

func orDone(done <-chan struct{}, ch <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case fName, ok := <-ch:
				if !ok {
					return
				}

				select {
				case <-done:
				case out <- fName:
				}
			}
		}

	}()

	return out
}

func generator(ctx context.Context, dir string) <-chan string {
	out := make(chan string, 128)
	go func() {
		defer close(out)
		r, _ := regexp.Compile(`(?i)\.(jpg|jpeg|png)$`)
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("readdir %s: %v\n", dir, err)
			return
		}
		for _, e := range entries {
			if e.IsDir() || !r.MatchString(e.Name()) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- e.Name():
			}
		}
	}()
	return out
}

func Decode(file *os.File, fileName string) (image.Image, error) {
	r1, _ := regexp.Compile(`(?i)\.(jpg|jpeg)$`)
	if r1.MatchString(fileName) {
		return jpeg.Decode(file)
	} else {
		return png.Decode(file)
	}
}
func Encode(file *os.File, m *image.Image) error {
	return jpeg.Encode(file, *m, nil)
}

func worker(OutputDir string, InputDir string, in <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for fName := range in {
		file, err := os.Open(filepath.Join(InputDir, fName))
		if err != nil {
			log.Fatal(err)
		}
		img, err := Decode(file, fName)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
		m := resize.Resize(100, 0, img, resize.Lanczos3)
		out, err := os.Create(filepath.Join(OutputDir, fName))
		if err != nil {
			log.Fatal(err)
		}
		Encode(out, &m)
	}

}

func GetThumbNail(ctx context.Context, in <-chan string, wg *sync.WaitGroup, OutputDir, InputDir string) {
	numWorkers := runtime.NumCPU()
	wg.Add(numWorkers)
	err := os.Mkdir(OutputDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	cin := orDone(ctx.Done(), in)

	for i := 0; i < numWorkers; i++ {
		go worker(OutputDir, InputDir, cin, wg)
	}
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(10 * time.Second)
		cancel()
	}()
	GetThumbNail(ctx, generator(ctx, `D:\6th sem\dbms`), &wg, `ThumbNails\`, `D:\6th sem\dbms`)
	wg.Wait()
	fmt.Println(time.Since(start))
}
