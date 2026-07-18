package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func jsonExercise() {
	_, _ = json.Marshal(&Person{Name: "test", Age: 30})
	_ = json.Unmarshal([]byte(`{"name":"test","age":30}`), &Person{})
}

func main() {
	jsonExercise()
	fmt.Println("json e2e test done")
}
