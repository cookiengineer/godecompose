package main

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func csvExercise() {
	r := csv.NewReader(strings.NewReader("a,b,c\n1,2,3\n"))
	records, _ := r.ReadAll()
	_ = records

	var buf strings.Builder
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"x", "y", "z"})
	w.Flush()
}

func main() {
	csvExercise()
	fmt.Println("csv e2e test done")
}
