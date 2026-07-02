package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	val  int
	next *Node
}

type LinkedList struct {
	head      *Node
	readCount int
	mutex     sync.Mutex
	turnstile sync.Mutex
	rwdLock   sync.Mutex
}

func (l *LinkedList) ReadLock() {
	l.turnstile.Lock()
	l.turnstile.Unlock()

	l.mutex.Lock()
	l.readCount++
	if l.readCount == 1 {
		l.rwdLock.Lock()
	}
	l.mutex.Unlock()
}

func (l *LinkedList) ReadUnlock() {
	l.mutex.Lock()
	l.readCount--
	if l.readCount == 0 {
		l.rwdLock.Unlock()
	}
	l.mutex.Unlock()
}

func (l *LinkedList) UpdateLock() {
	l.turnstile.Lock()
	l.rwdLock.Lock()
}

func (l *LinkedList) UpdateUnlock() {
	l.turnstile.Unlock()
	l.rwdLock.Unlock()
}

func (l *LinkedList) Contains(val int) bool {
	l.ReadLock()
	defer l.ReadUnlock()

	curr := l.head
	for curr != nil {
		if curr.val == val {
			return true
		}
		curr = curr.next
	}
	return false
}

func (l *LinkedList) Insert(val int) {
	l.UpdateLock()
	defer l.UpdateUnlock()

	newNode := &Node{val: val, next: l.head}
	l.head = newNode
}

func (l *LinkedList) Delete(val int) {
	l.UpdateLock()
	defer l.UpdateUnlock()

	if l.head == nil {
		return
	}
	if l.head.val == val {
		l.head = l.head.next
		return
	}

	curr := l.head
	for curr.next != nil {
		if curr.next.val == val {
			curr.next = curr.next.next
			return
		}
		curr = curr.next
	}
}

func main() {
	list := &LinkedList{}

	list.Insert(10)
	list.Insert(20)

	go func() {
		for {
			fmt.Printf("Search 10: %v\n", list.Contains(10))
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		list.Insert(30)
		fmt.Println("Inserted 30")
	}()

	go func() {
		time.Sleep(time.Second)
		list.Delete(10)
		fmt.Println("Deleted 10")
	}()

	time.Sleep(time.Second * 3)
}
