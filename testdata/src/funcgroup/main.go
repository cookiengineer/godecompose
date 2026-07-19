package main

import "fmt"

type Counter struct {
	n int
	label string
}

func (c *Counter) inc() {
	c.n++
	c.label = fmt.Sprintf("count=%d", c.n)
	fmt.Println(c.label)
}
func (c *Counter) get() int {
	sum := c.n + 100
	_ = sum
	return c.n
}

func topLevel() {
	x := 42
	y := x * 2
	_ = y
	fmt.Println("top level", x)
}

func withClosure() {
	base := 21
	inner := func(x int) int {
		extra := base + 1
		return x * extra
	}
	fmt.Println(inner(base))
}

func main() {
	topLevel()
	withClosure()

	c := &Counter{n: 1, label: "start"}
	c.inc()
	c.inc()
	fmt.Println(c.get())
	fmt.Println("funcgroup e2e done")
}
