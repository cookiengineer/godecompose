package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string { return "hello" }

func reflectExercise() {
	p := Person{Name: "test", Age: 30}

	t := reflect.TypeOf(p)
	_ = t.Name()
	_ = t.Kind()
	_ = t.NumMethod()
	_ = t.NumField()
	f := t.Field(0)
	_ = f
	_ = t.Size()
	_ = t.String()
	_ = t.Elem()

	_ = reflect.TypeOf(&p).Elem()

	v := reflect.ValueOf(p)
	_ = v.Interface()
	_ = v.Field(0).String()
	_ = v.Kind()
	_ = v.IsValid()
	_ = reflect.DeepEqual(p, Person{Name: "test", Age: 30})
}

func main() {
	reflectExercise()
	fmt.Println("reflect e2e test done")
}
