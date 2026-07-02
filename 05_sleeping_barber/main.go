package main

import (
	"fmt"
	"sync"
	"time"
)

const nChairs = 3

type Barbershop struct {
	waiting      int
	mutex        sync.Mutex
	customer     chan struct{}
	barber       chan struct{}
	haircutDone  chan struct{}
	customerLeft chan struct{}
}

func main() {
	shop := &Barbershop{
		customer:     make(chan struct{}, nChairs),
		barber:       make(chan struct{}),
		haircutDone:  make(chan struct{}),
		customerLeft: make(chan struct{}),
	}

	go shop.barberProcess()

	for i := 1; i <= 10; i++ {
		go shop.customerProcess(i)
		time.Sleep(time.Millisecond * 600)
	}

	time.Sleep(time.Second * 10)
}

func (s *Barbershop) barberProcess() {
	for {
		fmt.Println("Berberin spava...")
		<-s.customer

		s.barber <- struct{}{}
		fmt.Println("Berberin preuzima mušteriju i počinje šišanje.")

		time.Sleep(time.Second * 2)

		fmt.Println("Berberin je završio šišanje.")
		s.haircutDone <- struct{}{}
		<-s.customerLeft
	}
}

func (s *Barbershop) customerProcess(id int) {
	fmt.Printf("Mušterija %d dolazi u berbernicu.\n", id)

	s.mutex.Lock()
	if s.waiting == nChairs {
		fmt.Printf("Čekaonica puna! Mušterija %d odlazi bez šišanja.\n", id)
		s.mutex.Unlock()
		return
	}

	s.waiting++
	fmt.Printf("Mušterija %d seda u čekaonicu (Zauzeto: %d/%d).\n", id, s.waiting, nChairs)
	s.mutex.Unlock()

	s.customer <- struct{}{}
	<-s.barber

	s.mutex.Lock()
	s.waiting--
	s.mutex.Unlock()

	fmt.Printf("Mušterija %d se šiša...\n", id)
	<-s.haircutDone
	s.customerLeft <- struct{}{}
	fmt.Printf("Mušterija %d izlazi iz berbernice.\n", id)
}
