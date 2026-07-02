package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5

type Semaphore chan struct{}

func (s Semaphore) Wait()   { s <- struct{}{} }
func (s Semaphore) Signal() { <-s }

func main() {
	forks := make([]*sync.Mutex, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &sync.Mutex{}
	}

	footman := make(Semaphore, 4)

	var wg sync.WaitGroup

	for i := 0; i < numPhilosophers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			philosopher(id, forks[id], forks[(id+1)%numPhilosophers], footman)
		}(i)
	}

	wg.Wait()
}

func philosopher(id int, leftFork, rightFork *sync.Mutex, footman Semaphore) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Filozof %d razmišlja...\n", id)
		time.Sleep(time.Second)

		fmt.Printf("Filozof %d je ogladneo i traži dozvolu da sedne.\n", id)
		footman.Wait()

		leftFork.Lock()
		fmt.Printf("Filozof %d je uzeo LEVI štapić.\n", id)

		rightFork.Lock()
		fmt.Printf("Filozof %d je uzeo DESNI štapić i POČINJE DA JEDE.\n", id)

		time.Sleep(time.Second * 2)
		fmt.Printf("Filozof %d je završio jelo i vraća štapiće.\n", id)

		rightFork.Unlock()
		leftFork.Unlock()

		footman.Signal()
	}
}
