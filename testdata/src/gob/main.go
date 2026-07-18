package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Data struct {
	Name string
	Age  int
}

func gobExercise() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(&Data{Name: "test", Age: 30})

	dec := gob.NewDecoder(&buf)
	var d Data
	_ = dec.Decode(&d)

	gob.Register(&Data{})
	gob.RegisterName("mypkg.Data", &Data{})
}

func main() {
	gobExercise()
	fmt.Println("gob e2e test done")
}
