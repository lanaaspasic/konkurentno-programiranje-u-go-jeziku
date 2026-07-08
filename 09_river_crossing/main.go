package main

import (
	"fmt"
	"sync"
	"time"
)

type RiverCrossing struct {
	mutex       sync.Mutex
	hackers     int
	serfs       int
	hackerQueue chan struct{}
	serfQueue   chan struct{}
	barrier     sync.WaitGroup
}

func (rc *RiverCrossing) Hacker(id int) {
	isCaptain := false
	rc.mutex.Lock()
	rc.hackers++

	if rc.hackers == 4 {
		for i := 0; i < 4; i++ {
			rc.hackerQueue <- struct{}{}
		}
		rc.hackers = 0
		isCaptain = true
	} else if rc.hackers == 2 && rc.serfs >= 2 {
		for i := 0; i < 2; i++ {
			rc.hackerQueue <- struct{}{}
		}
		for i := 0; i < 2; i++ {
			rc.serfQueue <- struct{}{}
		}
		rc.serfs -= 2
		rc.hackers = 0
		isCaptain = true
	} else {
		rc.mutex.Unlock()
	}

	<-rc.hackerQueue
	fmt.Printf("Haker %d je ušao u čamac.\n", id)
	rc.barrier.Done()

	if isCaptain {
		rc.barrier.Wait()
		fmt.Println("--- KAPETAN (Haker): Svi su ukrcani. Krećemo! ---")
		rc.barrier.Add(4)
		rc.mutex.Unlock()
	}
}

func (rc *RiverCrossing) Serf(id int) {
	isCaptain := false
	rc.mutex.Lock()
	rc.serfs++

	if rc.serfs == 4 {
		for i := 0; i < 4; i++ {
			rc.serfQueue <- struct{}{}
		}
		rc.serfs = 0
		isCaptain = true
	} else if rc.serfs == 2 && rc.hackers >= 2 {
		rc.serfQueue <- struct{}{}
		rc.serfQueue <- struct{}{}
		rc.hackerQueue <- struct{}{}
		rc.hackerQueue <- struct{}{}

		rc.hackers -= 2
		rc.serfs = 0
		isCaptain = true
	} else {
		rc.mutex.Unlock()
	}

	<-rc.serfQueue
	fmt.Printf("Zaposleni %d je ušao u čamac.\n", id)
	rc.barrier.Done()

	if isCaptain {
		rc.barrier.Wait()
		fmt.Println("--- KAPETAN (Zaposleni): Svi su ukrcani. Krećemo! ---")
		rc.barrier.Add(4)
		rc.mutex.Unlock()
	}
}

func main() {
	rc := &RiverCrossing{
		hackerQueue: make(chan struct{}, 4),
		serfQueue:   make(chan struct{}, 4),
	}
	rc.barrier.Add(4)

	fmt.Println("Početak simulacije prelaska reke...")

	for i := 1; i <= 8; i++ {
		go rc.Hacker(i)
		go rc.Serf(i)
		time.Sleep(time.Millisecond * 100)
	}

	time.Sleep(time.Second * 10)
	fmt.Println("Simulacija završena.")
}
