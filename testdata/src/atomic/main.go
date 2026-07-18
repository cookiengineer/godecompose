package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

func atomicExercise() {
	var i32 int32
	atomic.StoreInt32(&i32, 42)
	_ = atomic.LoadInt32(&i32)
	atomic.AddInt32(&i32, 1)
	_ = atomic.SwapInt32(&i32, 0)
	atomic.CompareAndSwapInt32(&i32, 0, 100)

	var i64 int64
	atomic.StoreInt64(&i64, 42)
	_ = atomic.LoadInt64(&i64)
	atomic.AddInt64(&i64, 1)

	var u32 uint32
	atomic.StoreUint32(&u32, 42)
	_ = atomic.LoadUint32(&u32)

	var u64 uint64
	atomic.StoreUint64(&u64, 42)
	_ = atomic.LoadUint64(&u64)

	var up uintptr
	atomic.StoreUintptr(&up, 0x1234)
	_ = atomic.LoadUintptr(&up)

	var p unsafe.Pointer
	atomic.StorePointer(&p, nil)
	_ = atomic.LoadPointer(&p)
}

func main() {
	atomicExercise()
	fmt.Println("atomic e2e test done")
}
