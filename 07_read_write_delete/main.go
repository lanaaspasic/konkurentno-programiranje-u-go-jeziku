package main

import (
	"fmt"
	"sync"
	"time"
)

type Lightswitch struct {
	counter int
	mutex   sync.Mutex
}

func (l *Lightswitch) Lock(semaphore *sync.Mutex) {
	l.mutex.Lock()
	l.counter++
	if l.counter == 1 {
		semaphore.Lock()
	}
	l.mutex.Unlock()
}

func (l *Lightswitch) Unlock(semaphore *sync.Mutex) {
	l.mutex.Lock()
	l.counter--
	if l.counter == 0 {
		semaphore.Unlock()
	}
	l.mutex.Unlock()
}

type ReadWriteDelete struct {
	insertMutex  sync.Mutex
	noSearcher   sync.Mutex
	noInserter   sync.Mutex
	searchSwitch Lightswitch
	insertSwitch Lightswitch
}

func (rwd *ReadWriteDelete) Searcher(id int) {
	for i := 0; i < 3; i++ {
		rwd.searchSwitch.Lock(&rwd.noSearcher)

		fmt.Printf("[Searcher %d] Pretražuje podatke...\n", id)
		time.Sleep(time.Millisecond * 500)

		rwd.searchSwitch.Unlock(&rwd.noSearcher)
		time.Sleep(time.Second)
	}
}

func (rwd *ReadWriteDelete) Inserter(id int) {
	for i := 0; i < 3; i++ {
		rwd.insertSwitch.Lock(&rwd.noInserter)

		rwd.insertMutex.Lock()

		fmt.Printf("[Inserter %d] DODAJE nove podatke...\n", id)
		time.Sleep(time.Second)

		rwd.insertMutex.Unlock()
		rwd.insertSwitch.Unlock(&rwd.noInserter)

		time.Sleep(time.Second)
	}
}

func (rwd *ReadWriteDelete) Deleter(id int) {
	for i := 0; i < 3; i++ {
		rwd.noSearcher.Lock()
		rwd.noInserter.Lock()

		fmt.Printf("[!!! Deleter %d !!!] BRIŠE podatke - APSOLUTNI PRISTUP\n", id)
		time.Sleep(time.Second * 2)

		rwd.noInserter.Unlock()
		rwd.noSearcher.Unlock()

		time.Sleep(time.Second * 2)
	}
}

func main() {
	rwd := &ReadWriteDelete{}

	for i := 1; i <= 3; i++ {
		go rwd.Searcher(i)
	}
	for i := 1; i <= 2; i++ {
		go rwd.Inserter(i)
	}
	go rwd.Deleter(1)

	time.Sleep(time.Second * 15)
}
