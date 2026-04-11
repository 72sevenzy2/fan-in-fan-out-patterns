package main

import (
	"fmt"
	"sync"
)

/*
	exercise; fan-out -> generate 100 nums, create a number of workers to square those nums, and fan-in the results
	into one channel, and loop over the channel printing the results in it
*/

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func worker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// fanin (merge multiple channels into one)
func fanIn(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(chans))

	for _, ch := range chans {
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// generating the nums
	nums := make([]int, 0, 100)
	for i := 1; i <= 100; i++ {
		nums = append(nums, i)
	}

	in := generator(nums...)

	// send fan out workers
	w1 := worker(1, in)
	w2 := worker(2, in)
	w3 := worker(3, in)
	w4 := worker(4, in)
	w5 := worker(5, in)

	// merge into 1 fan in result
	out := fanIn(w1, w2, w3, w4, w5)

	// loop and print
	for v := range out {
		fmt.Println(v)
	}
}