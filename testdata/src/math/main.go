package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func mathExercise() {
	_ = math.Abs(-42.0)
	_ = math.Max(1.0, 2.0)
	_ = math.Min(1.0, 2.0)
	_ = math.Sqrt(16.0)
	_ = math.Pow(2.0, 3.0)
	_ = math.Sin(0.0)
	_ = math.Cos(0.0)
	_ = math.Tan(0.0)
	_ = math.Floor(3.14)
	_ = math.Ceil(3.14)
	_ = math.Round(3.14)
	_ = math.Log(1.0)
	_ = math.Log2(8.0)
	_ = math.Log10(100.0)
	_ = math.Exp(1.0)
	_ = math.Mod(10.0, 3.0)
	_ = math.Remainder(10.0, 3.0)
	_ = math.Hypot(3.0, 4.0)
}

func randExercise() {
	_ = rand.Intn(100)
	_ = rand.Float64()
	_ = rand.Int()
	_ = rand.Int31()
	_ = rand.Int63()
	_ = rand.Perm(10)
	rand.Shuffle(10, func(i, j int) {})
	rand.Seed(42)
	r := rand.New(rand.NewSource(42))
	_ = r
	_, _ = rand.Read(make([]byte, 4))
}

func sortExercise() {
	ints := []int{3, 1, 2}
	strings := []string{"c", "a", "b"}
	floats := []float64{3.0, 1.0, 2.0}

	sort.Ints(ints)
	sort.Strings(strings)
	sort.Float64s(floats)
	sort.Slice(ints, func(i, j int) bool { return ints[i] < ints[j] })
	_ = sort.Search(3, func(i int) bool { return ints[i] >= 2 })
	_ = sort.SearchInts(ints, 2)
	_ = sort.SearchStrings(strings, "b")
	_ = sort.SearchFloat64s(floats, 2.0)
	sort.Stable(sort.IntSlice(ints))
	_ = sort.IsSorted(sort.IntSlice(ints))
	_ = sort.Reverse(sort.IntSlice(ints))
}

func main() {
	mathExercise()
	randExercise()
	sortExercise()
	fmt.Println("math e2e test done")
}
