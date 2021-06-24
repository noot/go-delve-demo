package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func writeToCh(ch chan<- string, data string) {
	ch <- data
}

func writeToChA(ctx context.Context, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	str := "a"

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("writeToChA writing to ch")
			writeToCh(ch, str)
		}
	}
}

func writeToChB(ctx context.Context, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	str := "b"

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("writeToChB writing to ch")
			ch <- str
		}
	}
}

func readFromCh(ctx context.Context, wg *sync.WaitGroup, ch <-chan string) {
	defer wg.Done()
	
	for {
		select {
		case <-ctx.Done():
			return
		case str := <-ch:
			fmt.Println(str)
		}
	}
}

func main() {
	ch := make(chan string)
	ctx := context.Background()
	wg := new(sync.WaitGroup)

	go writeToChA(ctx, wg, ch)
	go writeToChB(ctx, wg, ch)
	go readFromCh(ctx, wg, ch)

	wg.Add(3)
	wg.Wait()
}