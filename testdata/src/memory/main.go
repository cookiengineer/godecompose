package main

import "fmt"

func memExercise() {
	src := make([]byte, 1024)
	dst := make([]byte, 1024)
	copy(dst, src)
	_ = dst
}

func appendExercise() {
	var data []int
	for i := 0; i < 10; i++ {
		data = append(data, i*2)
	}
	_ = data
}

func main() {
	memExercise()
	appendExercise()
	fmt.Println("memory e2e test done")
}
