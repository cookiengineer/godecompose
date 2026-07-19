package main

import "fmt"

func switchInt(v int) string {
	switch v {
	case 0:
		return "zero"
	case 1:
		return "one"
	case 2:
		return "two"
	default:
		return "many"
	}
}

func switchString(s string) string {
	switch s {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	default:
		return "other"
	}
}

func switchMultiCase(v int) string {
	switch v {
	case 0, 1, 2:
		return "small"
	case 3, 4:
		return "medium"
	default:
		return "large"
	}
}

func main() {
	fmt.Println(switchInt(1))
	fmt.Println(switchInt(5))
	fmt.Println(switchString("a"))
	fmt.Println(switchString("c"))
	fmt.Println(switchMultiCase(0))
	fmt.Println(switchMultiCase(3))
	fmt.Println("switchcase e2e done")
}
