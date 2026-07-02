package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	agentSem = make(chan struct{}, 1)
	tobacco  = make(chan struct{})
	paper    = make(chan struct{})
	match    = make(chan struct{})

	tobaccoSem = make(chan struct{})
	paperSem   = make(chan struct{})
	matchSem   = make(chan struct{})

	mutex                       sync.Mutex
	isTobacco, isPaper, isMatch bool
)

func main() {
	agentSem <- struct{}{}

	go agent()

	go pusherTobacco()
	go pusherPaper()
	go pusherMatch()

	go smokerWithTobacco()
	go smokerWithPaper()
	go smokerWithMatch()

	time.Sleep(time.Second * 15)
}

func agent() {
	for {
		<-agentSem
		t := time.Now().UnixNano() % 3
		switch t {
		case 0:
			fmt.Println("\n--- Agent stavlja: DUVAN i PAPIR ---")
			tobacco <- struct{}{}
			paper <- struct{}{}
		case 1:
			fmt.Println("\n--- Agent stavlja: PAPIR i ŠIBICE ---")
			paper <- struct{}{}
			match <- struct{}{}
		case 2:
			fmt.Println("\n--- Agent stavlja: DUVAN i ŠIBICE ---")
			tobacco <- struct{}{}
			match <- struct{}{}
		}
	}
}

func pusherTobacco() {
	for {
		<-tobacco
		mutex.Lock()
		if isPaper {
			isPaper = false
			matchSem <- struct{}{}
		} else if isMatch {
			isMatch = false
			paperSem <- struct{}{}
		} else {
			isTobacco = true
		}
		mutex.Unlock()
	}
}

func pusherPaper() {
	for {
		<-paper
		mutex.Lock()
		if isTobacco {
			isTobacco = false
			matchSem <- struct{}{}
		} else if isMatch {
			isMatch = false
			tobaccoSem <- struct{}{}
		} else {
			isPaper = true
		}
		mutex.Unlock()
	}
}

func pusherMatch() {
	for {
		<-match
		mutex.Lock()
		if isTobacco {
			isTobacco = false
			paperSem <- struct{}{}
		} else if isPaper {
			isPaper = false
			tobaccoSem <- struct{}{}
		} else {
			isMatch = true
		}
		mutex.Unlock()
	}
}

func smokerWithTobacco() {
	for {
		<-tobaccoSem
		fmt.Println("Pušač sa DUVANOM: Uzima papir i šibice, pravi cigaretu...")
		time.Sleep(time.Second)
		agentSem <- struct{}{}
		fmt.Println("Pušač sa DUVANOM: Puši...")
	}
}

func smokerWithPaper() {
	for {
		<-paperSem
		fmt.Println("Pušač sa PAPIROM: Uzima duvan i šibice, pravi cigaretu...")
		time.Sleep(time.Second)
		agentSem <- struct{}{}
		fmt.Println("Pušač sa PAPIROM: Puši...")
	}
}

func smokerWithMatch() {
	for {
		<-matchSem
		fmt.Println("Pušač sa ŠIBICAMA: Uzima duvan i papir, pravi cigaretu...")
		time.Sleep(time.Second)
		agentSem <- struct{}{}
		fmt.Println("Pušač sa ŠIBICAMA: Puši...")
	}
}
