package main

import "fmt"

func branch(x int) string {
	if x > 0 {
		return "positive"
	} else if x < 0 {
		return "negative"
	}
	return "zero"
}

func loopCount(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

func main() {
	r1 := branch(10)
	r2 := branch(-5)
	r3 := branch(0)
	s := loopCount(5)
	fmt.Println(r1, r2, r3, s)
	fmt.Println("structout e2e done")
}
