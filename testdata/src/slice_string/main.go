package main

import "fmt"

func sliceStringExercise() {
	nums := make([]int, 5, 10)
	_ = nums

	a := []byte("hello")
	s := string(a)
	_ = s

	b := "world"
	bs := []byte(b)
	_ = bs

	greeting := "hello" + " " + "world"
	_ = greeting
}

func main() {
	sliceStringExercise()
	fmt.Println("slice/string e2e test done")
}
