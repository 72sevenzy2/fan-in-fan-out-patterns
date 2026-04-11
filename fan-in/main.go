package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan<- string) {
	for range 3 {
		ch <- fmt.Sprintf("received work with id %d", id)
		time.Sleep(time.Millisecond * 500)
	}
	close(ch)
}

func fanin(channels ...<-chan string) <-chan string {
	out := make(chan string)
	for _, value := range channels {
		go func (val <-chan string)  {
			for i := range val {
				out <- i
			}
		}(value)
	}

	go func ()  {
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

	merged := fanin(ch1, ch2)

	for val := range merged {
		fmt.Println(val)
	}
}