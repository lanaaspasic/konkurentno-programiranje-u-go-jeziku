package main

import (
	"fmt"
	"time"
)

func main() {
	semaphore := make(chan struct{}, 3)

	for i := 1; i <= 10; i++ {
		go func(id int) {
			fmt.Printf("Radnik %d čeka na slobodno mesto...\n", id)

			semaphore <- struct{}{}

			fmt.Printf("--- Radnik %d je ušao i radi ---\n", id)
			time.Sleep(2 * time.Second)

			<-semaphore

			fmt.Printf("Radnik %d je završio.\n", id)
		}(i)
	}

	time.Sleep(10 * time.Second)
}
