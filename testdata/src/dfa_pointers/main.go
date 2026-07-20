package main

import "fmt"

type Item struct {
	Name  string
	Value int
}

func updateItem(it *Item, v int) {
	it.Value = v
}

func printItem(it *Item) {
	fmt.Println(it.Name, it.Value)
}

func main() {
	it := &Item{Name: "test", Value: 42}
	updateItem(it, 100)
	printItem(it)
}
