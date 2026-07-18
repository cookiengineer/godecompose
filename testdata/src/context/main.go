package main

import (
	"context"
	"fmt"
	"time"
)

func contextExercise() {
	ctx := context.Background()
	_ = ctx

	todo := context.TODO()
	_ = todo

	ctx, cancel := context.WithCancel(todo)
	defer cancel()

	ctxTO, cancelTO := context.WithTimeout(todo, time.Second)
	defer cancelTO()
	_ = ctxTO

	ctxDL, cancelDL := context.WithDeadline(todo, time.Now().Add(time.Second))
	defer cancelDL()
	_ = ctxDL

	ctxVal := context.WithValue(todo, "key", "value")
	_ = ctxVal.Value("key")
}

func main() {
	contextExercise()
	fmt.Println("context e2e test done")
}
