package main

import (
	"fmt"
	"os"
)

func ifNilCheck(x *int) {
	if x != nil {
		fmt.Println(*x)
	}
}

func ifErrNotNil() error {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func ifElseComparison(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ifNotEqualConst(v int) {
	if v != 0 {
		fmt.Println(v)
	}
}

func ifLessThan(v int) {
	if v < 10 {
		fmt.Println("small")
	}
}

func ifNotNilReturnNilErr() error {
	_, err := os.Stat("/tmp")
	if err != nil {
		return err
	}
	return nil
}

func deferFunc() {
	f, err := os.CreateTemp("", "defer")
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Println("deferred")
}

func goroutineFunc() {
	ch := make(chan int, 1)
	go func() {
		ch <- 42
	}()
	val := <-ch
	fmt.Println(val)
}

func mapOperations() {
	m := make(map[string]int)
	m["key"] = 1
	v := m["key"]
	_ = v
	v2, ok := m["missing"]
	if !ok {
		fmt.Println("not found:", v2)
	}
}

func channelOps() {
	ch := make(chan int, 1)
	ch <- 1
	val := <-ch
	fmt.Println(val)
}

func panicOnError() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	_, err := os.Stat("/nonexistent")
	if err != nil {
		panic(err)
	}
}

func main() {
	n := 42
	ifNilCheck(&n)

	err := ifErrNotNil()
	_ = err

	max := ifElseComparison(5, 10)
	_ = max

	ifNotEqualConst(7)
	ifLessThan(3)

	_ = ifNotNilReturnNilErr()

	deferFunc()
	goroutineFunc()
	mapOperations()
	channelOps()
	panicOnError()

	fmt.Println("controlflow e2e done")
}
