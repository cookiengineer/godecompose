package main

import (
	"fmt"
	"sync"
)

func syncExercise() {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()

	var once sync.Once
	once.Do(func() {
		fmt.Println("ran once")
	})
}

func main() {
	syncExercise()
}
