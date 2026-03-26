package main

import (
	"fmt"
	"sync"
	"time"
)

// Counter is a simple counter with mutex protection
type Counter struct {
	mu    sync.Mutex
	value int
}

// Increment increments the counter
func (c *Counter) Increment() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// Value returns the current value
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Worker demonstrates a goroutine with channels
func Worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		// Simulate work
		time.Sleep(10 * time.Millisecond)
		results <- j * 2
	}
}

// Generator returns a channel that produces numbers
func Generator(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

// SelectExample demonstrates select with multiple channels
func SelectExample(ch1, ch2 <-chan string, timeout time.Duration) string {
	select {
	case msg := <-ch1:
		return msg
	case msg := <-ch2:
		return msg
	case <-time.After(timeout):
		return "timeout"
	}
}

func main() {
	counter := &Counter{}
	fmt.Printf("Initial value: %d\n", counter.Value())

	counter.Increment()
	fmt.Printf("After increment: %d\n", counter.Value())

	// Example of using channels
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for i := 1; i <= 3; i++ {
		go Worker(i, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for r := 1; r <= 5; r++ {
		fmt.Printf("Result: %d\n", <-results)
	}

	// Example of panic/recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// This would panic but is caught by recover above
	// panic("test panic")
}
