package main

import (
	"fmt"
	"time"
	"unsafe"
)

type Point struct {
	X int
	Y int
}

func intAssignment() {
	var a int
	a = 42
	_ = a

	b := 100
	_ = b

	c := a + b
	fmt.Println(c)
}

func stringCreation() {
	s := "hello world"
	_ = s

	b := []byte(s)
	_ = b

	b2 := []byte("test")
	s2 := string(b2)
	_ = s2
}

func boolOps() {
	t := true
	f := false
	if t && !f {
		fmt.Println("bool ops")
	}
}

func sliceOps() {
	sl := make([]int, 10)
	sl[0] = 1
	sl[1] = 2

	sl2 := make([]string, 3)
	sl2[0] = "a"
	_ = sl2

	sl = append(sl, 3, 4, 5)
	fmt.Println(len(sl))

	dst := make([]int, len(sl))
	copy(dst, sl)
	fmt.Println(dst[0])

	sl3 := []int{1, 2, 3}
	_ = sl3
}

func mapOps() {
	m := make(map[string]int)
	m["one"] = 1
	m["two"] = 2

	v := m["one"]
	_ = v

	delete(m, "two")
}

func structOps() {
	p := Point{X: 10, Y: 20}
	_ = p

	p2 := new(Point)
	p2.X = 30
	p2.Y = 40
	fmt.Println(p2.X, p2.Y)

	p3 := &Point{X: 50}
	fmt.Println(p3.Y)
}

func pointerOps() {
	n := 42
	p := &n
	*p = 100
	fmt.Println(*p)
}

func uintptrOps() {
	var x int = 42
	u := uintptr(unsafe.Pointer(&x))
	_ = u
}

func interfaceOps() {
	var v interface{}
	v = 42
	_ = v

	v = "hello"
	_ = v

	var v2 interface{} = []int{1, 2}
	_ = v2
}

func timeOps() {
	t := time.Now()
	_ = t
}

func concatOps() {
	a := "hello"
	b := "world"
	c := a + " " + b
	fmt.Println(c)
}

func clearOps() {
	buf := make([]byte, 100)
	for i := range buf {
		buf[i] = 0
	}

	clear(buf)
	fmt.Println(buf[0])
}

func main() {
	intAssignment()
	stringCreation()
	boolOps()
	sliceOps()
	mapOps()
	structOps()
	pointerOps()
	uintptrOps()
	interfaceOps()
	timeOps()
	concatOps()
	clearOps()

	fmt.Println("dataops e2e done")
}
