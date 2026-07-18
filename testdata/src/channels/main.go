package main

import "fmt"

func channelExercise() {
	ch := make(chan int, 1)
	ch <- 42
	val := <-ch
	close(ch)
	_ = val

	selectCh := make(chan string, 1)
	select {
	case selectCh <- "hello":
		<-selectCh
	default:
	}
}

func main() {
	channelExercise()
	fmt.Println("channels e2e test done")
}
