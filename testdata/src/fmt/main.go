package main

import "fmt"

func fmtExercise() {
	s := fmt.Sprintf("hello %s", "world")
	fmt.Println(s)
	fmt.Printf("value: %d\n", 42)
	_ = fmt.Errorf("test error: %d", 1)
}

func main() {
	fmtExercise()
}
