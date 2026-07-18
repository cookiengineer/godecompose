package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"
)

func multipartExercise() {
	body := "--boundary\r\nContent-Disposition: form-data; name=\"field\"\r\n\r\nvalue\r\n--boundary--\r\n"
	r := multipart.NewReader(strings.NewReader(body), "boundary")
	part, _ := r.NextPart()
	if part != nil {
		io.ReadAll(part)
	}

	var buf strings.Builder
	w := multipart.NewWriter(&buf)
	_ = w.FormDataContentType()
	w.Close()
}

func main() {
	multipartExercise()
	fmt.Println("multipart e2e test done")
}
