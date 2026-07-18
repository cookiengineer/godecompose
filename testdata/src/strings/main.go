package main

import (
	"fmt"
	"strings"
)

func stringsExercise() {
	s := "hello world"
	_ = strings.Contains(s, "world")
	_ = strings.HasPrefix(s, "hello")
	_ = strings.HasSuffix(s, "world")
	_ = strings.Index(s, "w")
	_ = strings.LastIndex(s, "o")
	_ = strings.Count(s, "l")
	_ = strings.Join([]string{"a", "b"}, ",")
	_ = strings.Split(s, " ")
	_ = strings.ReplaceAll(s, "world", "gophers")
	_ = strings.ToLower("HELLO")
	_ = strings.ToUpper("hello")
	_ = strings.TrimSpace("  hello  ")
	_ = strings.TrimPrefix(s, "hello")
	_ = strings.TrimSuffix(s, "world")
	_ = strings.Repeat("ha", 3)
	_ = strings.Fields("a b c")

	_ = strings.NewReader(s)
	_ = strings.NewReplacer("a", "b")

	var b strings.Builder
	b.WriteString("hello")
	b.WriteByte('!')
	b.WriteRune('世')
	_ = b.String()
	b.Reset()
}

func main() {
	stringsExercise()
	fmt.Println("strings e2e test done")
}
