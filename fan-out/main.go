package main

import (
	"fmt"
	"sync"
)

func worker(wg *sync.WaitGroup, id int, tasks <-chan int, result chan<- int) {
	defer wg.Done() // the wg decrements by 1 (defer makes it run no matter what after in this code block)
	for i := range tasks {
		fmt.Println("work received with", id)
		result <- i
	}
}

func main() {
	tasks := make(chan int) // jobs
	res := make(chan int) // results chan

	var wg sync.WaitGroup // initialising in the waitgroup (to synchronise the res channel closing)

	for i := range 3 {
		wg.Add(1) // adds 1 for each worker
		go worker(&wg, i, tasks, res) // seperate thread(gorountine) for each worker
	}
	
	go func() {
		for i := range 8 { // starting 8 jobs
			tasks <- i // sending tasks 
		}
		close(tasks) // closing the tasks afterwards (if not closed the workers will continue to read from tasks and when theres no more value to receive it causes a block.)
	}()

	go func() { 
		wg.Wait() // waits tills wg.Done() in worker func so that it can close the results channel
		close(res) // closing results as we're looping over them when receiving the res channel below, without this it would simply print all available values and wait forever. (block)
	}()

	for range res { 
		fmt.Println(<-res) // receiving and printing the values
	}
}
