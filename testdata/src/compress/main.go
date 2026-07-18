package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
)

func compressExercise() {
	var buf bytes.Buffer

	gw, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	_, _ = gw.Write([]byte("hello"))
	gw.Flush()
	gw.Close()

	_ = gzip.NewWriter(&buf)
	gr, _ := gzip.NewReader(&buf)
	if gr != nil {
		b := make([]byte, 64)
		_, _ = gr.Read(b)
		gr.Close()
	}

	_ = zlib.NewWriter(&buf)
	zw, _ := zlib.NewWriterLevel(&buf, zlib.DefaultCompression)
	_ = zw
	zr, _ := zlib.NewReader(&buf)
	_ = zr

	_ = flate.NewReader(&buf)
	fw, _ := flate.NewWriter(&buf, flate.DefaultCompression)
	fw.Close()
}

func main() {
	compressExercise()
	fmt.Println("compress e2e test done")
}
