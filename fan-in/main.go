package main

import (
	"fmt"
	"sync"
)

func worker(id int, ch chan<- string) {
	for range 3 {
		ch <- fmt.Sprintf("work received with id %d", id)
	}
	close(ch) // after assigning the ch with the work message 3 times, we close it
}

func fanin(wg *sync.WaitGroup, ch ...<-chan string) <-chan string { // returns a channel which is to be received
	out := make(chan string) // receiver channel
	wg.Add(len(ch)) // add length of the channels (if theres more than 1)

	for _, value := range ch { // ignore index and get each value of range ...ch
		go func(ch <-chan string) { // another thread which takes in a parameter which is value of the range ...ch we did in the outer loop.
			defer wg.Done() // decrements the wg by 1.
			for i := range ch { // looping over each channel
				out <- i // assign the receiver channel to the parameter ch
			}
		}(value) // calling the func immediately and passing the param as value (value of range ch in outer loop)
	}

	go func() { // seperate thread to synchronise out chan (receiver channel) to close
		wg.Wait() // waits till the wg is 0
		close(out)
	}()

	return out
}

func main() {
	ch1 := make(chan string) 
	ch2 := make(chan string)

	var wg sync.WaitGroup

	// seperate thread for both workers
	go worker(1, ch1)
	go worker(2, ch2)

	merged := fanin(&wg, ch1, ch2) // grabbing the merged results and printing them below

	for val := range merged {
		fmt.Println(val)
	}
}
