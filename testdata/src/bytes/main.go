package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func bytesExercise() {
	a := []byte("hello world")
	b := []byte("hello")

	_ = bytes.Contains(a, b)
	_ = bytes.Equal(a, b)
	_ = bytes.Compare(a, b)
	_ = bytes.Index(a, b)
	_ = bytes.LastIndex(a, b)
	_ = bytes.Join([][]byte{a, b}, []byte(","))
	_ = bytes.Split(a, b)
	_ = bytes.HasPrefix(a, b)

	_ = bytes.NewReader(a)

	buf := bytes.NewBuffer(a)
	buf.WriteString("hello")
	buf.WriteByte('!')
	_ = buf.Bytes()
	_ = buf.String()
	buf.Reset()
}

func bufioExercise() {
	s := "line1\nline2\nline3\n"

	_ = bufio.NewReader(strings.NewReader(s))
	_ = bufio.NewWriter(new(bytes.Buffer))
	_ = bufio.NewScanner(strings.NewReader(s))
}

func main() {
	bytesExercise()
	bufioExercise()
	fmt.Println("bytes e2e test done")
}
