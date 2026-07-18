package main

import (
	"fmt"
	"unsafe"
)

type S struct {
	a int32
	b int64
}

func unsafeExercise() {
	s := S{a: 1, b: 2}
	_ = unsafe.Sizeof(s)
	_ = unsafe.Offsetof(s.b)
	_ = unsafe.Alignof(s)

	p := unsafe.Pointer(&s)
	_ = p

	slice := []int{1, 2, 3}
	_ = unsafe.SliceData(slice)
	data := make([]byte, 16)
	_ = unsafe.Slice(&data[0], 16)
	_ = unsafe.Add(unsafe.Pointer(&data[0]), 8)

	str := "hello"
	_ = unsafe.StringData(str)
	_ = unsafe.String(unsafe.StringData(str), len(str))
}

func main() {
	unsafeExercise()
	fmt.Println("unsafe e2e test done")
}
