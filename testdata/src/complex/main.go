package main

import (
	"fmt"
	"sync"
)

type Calculator interface {
	Compute(a, b float64) float64
}

type Adder struct{}
func (Adder) Compute(a, b float64) float64 { return a + b }

type Multiplier struct{}
func (Multiplier) Compute(a, b float64) float64 { return a * b }

func processSlice(data []int) []int {
	result := make([]int, len(data))
	for i, v := range data {
		result[i] = v * 2
	}
	return result
}

func processMap(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func processChannel(ch chan int, n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += <-ch
	}
	return sum
}

func goroutineFanIn(values []int) int {
	ch := make(chan int, len(values))
	var wg sync.WaitGroup

	for _, v := range values {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			ch <- val * val
		}(v)
	}

	wg.Wait()
	close(ch)

	result := 0
	for v := range ch {
		result += v
	}
	return result
}

func useInterface(calc Calculator, a, b float64) float64 {
	return calc.Compute(a, b)
}

func main() {
	slice := processSlice([]int{1, 2, 3, 4, 5})
	keys := processMap(map[string]int{"a": 1, "b": 2, "c": 3})
	_ = slice
	_ = keys

	ch := make(chan int, 3)
	go func() {
		ch <- 10
		ch <- 20
		ch <- 30
	}()
	sum := processChannel(ch, 3)
	fmt.Printf("channel sum: %d\n", sum)

	fanSum := goroutineFanIn([]int{1, 2, 3, 4})
	fmt.Printf("fan-in sum: %d\n", fanSum)

	addResult := useInterface(Adder{}, 10.5, 20.5)
	mulResult := useInterface(Multiplier{}, 3.0, 4.0)
	fmt.Printf("add: %.1f, mul: %.1f\n", addResult, mulResult)
}
