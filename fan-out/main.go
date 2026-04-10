package main

import (
	"fmt"
	"sync"
)

func worker(wg *sync.WaitGroup, id int, tasks <-chan int, result chan<- int) {
	defer wg.Done()
	for i := range tasks {
		fmt.Println("work received with", id)
		result <- i
	}
}

func main() {
	tasks := make(chan int)
	res := make(chan int)

	var wg sync.WaitGroup

	for i := range 3 {
		wg.Add(1)
		go worker(&wg, i, tasks, res)
	}

	go func() {
		for i := range 8 {
			tasks <- i
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(res)
	}()

	for range res {
		fmt.Println(<-res)
	}
}
