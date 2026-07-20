package main

import "fmt"

func add42(x int) int {
	y := x + 42
	return y
}

func main() {
	fmt.Println(add42(10))
}
