package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func main() {
	// 1. Kreiranje test stabla
	//        1
	//       / \
	//      2   3
	//     / \
	//    4   5
	root := &Node{Value: 1,
		Left:  &Node{Value: 2, Left: &Node{Value: 4}, Right: &Node{Value: 5}},
		Right: &Node{Value: 3},
	}

	workQueue := make(chan *Node, 10)

	var wg sync.WaitGroup

	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		go func(workerID int) {
			for node := range workQueue {
				fmt.Printf("[Worker %d] Obrađuje čvor: %d\n", workerID, node.Value)
				time.Sleep(time.Millisecond * 500)

				if node.Left != nil {
					wg.Add(1)
					workQueue <- node.Left
				}
				if node.Right != nil {
					wg.Add(1)
					workQueue <- node.Right
				}

				wg.Done()
			}
		}(i)
	}

	wg.Add(1)
	workQueue <- root

	wg.Wait()
	close(workQueue)
	fmt.Println("Obilazak stabla završen.")
}
