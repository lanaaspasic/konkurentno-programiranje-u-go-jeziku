package main

import (
	"fmt"
	"sync"
	"time"
)

type Barrier struct {
	mutex sync.Mutex
	cond  *sync.Cond
	count int
	n     int
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	b.count++
	if b.count == b.n {
		b.count = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.mutex.Unlock()
}

type H2O struct {
	mutex      sync.Mutex
	hydrogen   int
	oxygen     int
	hydroQueue chan struct{}
	oxyQueue   chan struct{}
	barrier    *Barrier
}

func (h *H2O) Hydrogen(id int) {
	h.mutex.Lock()
	h.hydrogen++

	if h.hydrogen >= 2 && h.oxygen >= 1 {
		h.hydroQueue <- struct{}{}
		h.hydroQueue <- struct{}{}
		h.hydrogen -= 2
		h.oxyQueue <- struct{}{}
		h.oxygen -= 1
	} else {
		h.mutex.Unlock()
	}

	<-h.hydroQueue
	fmt.Printf("Vodonik %d: bond()\n", id)

	h.barrier.Wait()
	time.Sleep(time.Millisecond * 500)
}

func (h *H2O) Oxygen(id int) {
	h.mutex.Lock()
	h.oxygen++

	if h.hydrogen >= 2 {
		h.hydroQueue <- struct{}{}
		h.hydroQueue <- struct{}{}
		h.hydrogen -= 2
		h.oxyQueue <- struct{}{}
		h.oxygen -= 1
	} else {
		h.mutex.Unlock()
	}

	<-h.oxyQueue
	fmt.Printf("KISEONIK %d: bond() !!!\n", id)

	h.barrier.Wait()
	time.Sleep(time.Millisecond * 500)

	fmt.Println("--- Molekul H2O je uspešno formiran ---")
	time.Sleep(time.Second * 1)
	h.mutex.Unlock()
}

func main() {
	h2o := &H2O{
		hydroQueue: make(chan struct{}, 100),
		oxyQueue:   make(chan struct{}, 100),
		barrier:    NewBarrier(3),
	}

	fmt.Println("Početak simulacije formiranja molekula vode...")

	for i := 1; i <= 4; i++ {
		go h2o.Oxygen(i)
	}
	for i := 1; i <= 8; i++ {
		go h2o.Hydrogen(i)
	}

	time.Sleep(time.Second * 5)
	fmt.Println("Simulacija završena.")
}
