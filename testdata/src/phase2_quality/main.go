package main

import (
	"errors"
	"fmt"
)

type Counter struct {
	Value int
	Label string
}

func NewCounter(label string, start int) *Counter {
	return &Counter{Value: start, Label: label}
}

func (c *Counter) Increment() int {
	c.Value++
	return c.Value
}

func (c *Counter) GetValue() int {
	return c.Value
}

func (c *Counter) SetValue(v int) {
	c.Value = v
}

func processData(data []int, threshold int) (int, error) {
	if data == nil {
		return 0, errors.New("data is nil")
	}
	sum := 0
	count := 0
	for i, v := range data {
		if v > threshold {
			sum += v
			count++
		}
		_ = i
	}
	if count == 0 {
		return 0, nil
	}
	return sum / count, nil
}

func classify(n int) string {
	if n < 0 {
		return "negative"
	} else if n == 0 {
		return "zero"
	} else if n < 100 {
		return "small"
	}
	return "large"
}

func switchExample(op string, a, b int) int {
	switch op {
	case "add":
		return a + b
	case "sub":
		return a - b
	case "mul":
		return a * b
	case "div":
		if b != 0 {
			return a / b
		}
		return 0
	}
	return -1
}

func main() {
	c := NewCounter("items", 0)
	c.Increment()
	c.Increment()
	c.SetValue(5)
	result, err := processData([]int{1, 5, 10, 15}, 7)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(c.GetValue(), result)
	fmt.Println(classify(42))
	fmt.Println(switchExample("mul", 6, 7))
}
