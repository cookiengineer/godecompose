package main

import (
	"fmt"
	"strconv"
)

func strconvExercise() {
	_ = strconv.Itoa(42)
	_, _ = strconv.Atoi("42")
	_ = strconv.FormatInt(255, 16)
	_, _ = strconv.ParseInt("ff", 16, 64)
	_ = strconv.FormatFloat(3.14, 'f', 2, 64)
	_, _ = strconv.ParseFloat("3.14", 64)
	_ = strconv.Quote("hello")
	_, _ = strconv.Unquote(`"hello"`)
}

func main() {
	strconvExercise()
	fmt.Println("strconv e2e test done")
}
