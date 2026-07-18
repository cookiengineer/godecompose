package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
)

func archiveExercise() {
	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)
	_ = tw.WriteHeader(&tar.Header{Name: "test.txt", Size: 5})
	_, _ = tw.Write([]byte("hello"))
	tw.Flush()
	tw.Close()

	tr := tar.NewReader(&buf)
	hdr, _ := tr.Next()
	if hdr != nil {
		b := make([]byte, 64)
		_, _ = tr.Read(b)
	}

	var zbuf bytes.Buffer
	zw := zip.NewWriter(&zbuf)
	w, _ := zw.Create("test.txt")
	_, _ = w.Write([]byte("hello"))
	zw.Close()

	zr, _ := zip.NewReader(bytes.NewReader(zbuf.Bytes()), int64(zbuf.Len()))
	for _, f := range zr.File {
		rc, _ := f.Open()
		rc.Close()
	}
}

func main() {
	archiveExercise()
	fmt.Println("archive e2e test done")
}
