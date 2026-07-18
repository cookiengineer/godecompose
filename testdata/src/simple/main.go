package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func multiply(x, y int64) int64 {
	return x * y
}

func factorial(n uint64) uint64 {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	sum := add(3, 4)
	prod := multiply(5, 6)
	fact := factorial(7)

	fmt.Printf("add(3,4)=%d multiply(5,6)=%d factorial(7)=%d\n", sum, prod, fact)
}
