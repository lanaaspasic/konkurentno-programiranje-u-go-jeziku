package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const fileName = "shared.txt"

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

type ReaderWriterNoStarve struct {
	roomEmpty  sync.Mutex
	turnstile  sync.Mutex
	readSwitch Lightswitch
}

func (rw *ReaderWriterNoStarve) Reader(id int) {
	for i := 0; i < 3; i++ {

		rw.turnstile.Lock()
		rw.turnstile.Unlock()

		rw.readSwitch.Lock(&rw.roomEmpty)

		data, _ := os.ReadFile(fileName)
		fmt.Printf("[Čitalac %d - NoStarve] Čita: %s\n", id, string(data))
		time.Sleep(time.Millisecond * 500)

		rw.readSwitch.Unlock(&rw.roomEmpty)
		time.Sleep(time.Millisecond * 100)
	}
}

func (rw *ReaderWriterNoStarve) Writer(id int) {
	for i := 0; i < 3; i++ {
		rw.turnstile.Lock()
		rw.roomEmpty.Lock()

		f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		poruka := fmt.Sprintf("Pisac %d (No-Starve) piše %d\n", id, i)
		f.WriteString(poruka)
		fmt.Printf("!!! %s", poruka)
		f.Close()
		time.Sleep(time.Second)

		rw.turnstile.Unlock()
		rw.roomEmpty.Unlock()
		time.Sleep(time.Second)
	}
}

func main() {
	os.Create(fileName)

	rw := &ReaderWriterNoStarve{}
	for i := 1; i <= 5; i++ {
		go rw.Reader(i)
	}
	for i := 1; i <= 2; i++ {
		go rw.Writer(i)
	}

	time.Sleep(time.Second * 15)
}
