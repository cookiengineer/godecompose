package main

import "fmt"

func mapExercise() {
	m := make(map[string]int)
	m["alpha"] = 1
	m["beta"] = 2

	if val, ok := m["alpha"]; ok {
		_ = val
	}

	delete(m, "beta")

	for k := range m {
		_ = k
	}
}

func main() {
	mapExercise()
	fmt.Println("maps e2e test done")
}
