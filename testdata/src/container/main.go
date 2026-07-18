package main

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any          { n := len(*h); x := (*h)[n-1]; *h = (*h)[:n-1]; return x }

func containerExercise() {
	l := list.New()
	l.Init()
	l.PushFront(1)
	l.PushBack(2)
	_ = l.Front()
	_ = l.Back()
	_ = l.Len()
	e := l.Front()
	l.InsertBefore(3, e)
	l.InsertAfter(4, e)
	l.MoveToFront(e)
	l.MoveToBack(e)
	l.Remove(e)

	h := &IntHeap{2, 1, 3}
	heap.Init(h)
	heap.Push(h, 0)
	_ = heap.Pop(h)
	heap.Fix(h, 0)
	heap.Remove(h, 0)

	r := ring.New(3)
	_ = r.Len()
	_ = r.Next()
	_ = r.Prev()
	_ = r.Move(1)
	r.Do(func(x any) {})
	s := ring.New(1)
	_ = r.Link(s)
}

func main() {
	containerExercise()
	fmt.Println("container e2e test done")
}
