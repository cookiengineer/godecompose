package main

import (
	"fmt"
	"os"
)

func osExercise() {
	f, err := os.Create("/tmp/godecompose-e2e-test")
	if err != nil {
		fmt.Println("create error:", err)
		return
	}

	_, _ = f.Write([]byte("test data"))
	_ = f.Close()

	data, err := os.ReadFile("/tmp/godecompose-e2e-test")
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	_ = data

	_ = os.Remove("/tmp/godecompose-e2e-test")
}

func main() {
	osExercise()
	fmt.Println("os e2e test done")
}
