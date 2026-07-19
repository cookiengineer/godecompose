package main

import (
	"fmt"
	"os"
)

func countedLoop(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
		if sum > 100 {
			sum = 100
			break
		}
		if sum < 0 {
			sum = 0
		}
	}
	fmt.Println("sum:", sum)
	return sum
}

func errorReturn() error {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		fmt.Println("create fail")
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	buf := []byte("hello")
	_, err = f.Write(buf)
	if err != nil {
		fmt.Println("write fail")
		return fmt.Errorf("write: %w", err)
	}
	fmt.Println("write ok")
	return nil
}

func errorReturnNil() error {
	f, err := os.CreateTemp("", "test")
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	fmt.Println(f.Name())
	return nil
}

func typeAssert(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return "not a string"
	}
	return s
}

func main() {
	fmt.Println(countedLoop(5))
	if err := errorReturn(); err != nil {
		fmt.Println("err:", err)
	}
	if err := errorReturnNil(); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(typeAssert("hello"))
	fmt.Println(typeAssert(42))
	fmt.Println("looperror e2e done")
}
