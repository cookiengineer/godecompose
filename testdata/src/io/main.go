package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func ioExercise() {
	r := strings.NewReader("hello world")
	_, _ = io.ReadAll(r)
}

func timeExercise() {
	_ = time.Now()
	time.Sleep(0)
	_ = time.After(0)
	_ = time.Since(time.Now())
	_ = time.NewTicker(time.Second)
}

func main() {
	ioExercise()
	timeExercise()
	fmt.Println("io e2e test done")
}
