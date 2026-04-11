package main

import (
	"fmt"
	"sync"
)

func worker(id int, ch chan<- string) {
	for range 3 {
		ch <- fmt.Sprintf("work received with id %d", id)
	}
	close(ch)
}

func fanin(wg *sync.WaitGroup, ch ...<-chan string) <-chan string {
	out := make(chan string)
	wg.Add(len(ch))
	
	for _, value := range ch {
		go func (ch <-chan string)  {
			defer wg.Done()
			for i := range ch {
				out <- i
			}
		}(value)
	}

	go func ()  {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	var wg sync.WaitGroup

	go worker(1, ch1)
	go worker(2, ch2)

	merged := fanin(&wg, ch1, ch2)

	for val := range merged {
		fmt.Println(val)
	}
}