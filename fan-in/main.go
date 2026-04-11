package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan<- string) {
	for i := 0; i < 3; i++ {
		ch <- fmt.Sprintf("Worker %d: job %d", id, i)
		time.Sleep(time.Millisecond * 500)
	}
	close(ch)
}

func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	for _, ch := range channels {
		go func(c <-chan string) {
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		time.Sleep(2 * time.Second)
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go worker(1, ch1)
	go worker(2, ch2)

	merged := fanIn(ch1, ch2)

	for msg := range merged {
		fmt.Println(msg)
	}
}