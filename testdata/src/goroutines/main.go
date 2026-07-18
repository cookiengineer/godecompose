package main

import "fmt"

func goroutineExercise() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		ch <- 100
	}()

	val := <-ch
	_ = val
}

func main() {
	goroutineExercise()
	fmt.Println("goroutines e2e test done")
}
