package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	v   int
	mux sync.Mutex
}

func (c *SafeCounter) Inc() {
	c.mux.Lock()
	c.v++
	c.mux.Unlock()
}

func main() {
	c := SafeCounter{v: 0}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("Konačna vrednost brojača: %d (Očekivano: 1000)\n", c.v)
}
