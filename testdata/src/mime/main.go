package main

import (
	"fmt"
	"mime"
)

func mimeExercise() {
	_ = mime.TypeByExtension(".html")
	_, _ = mime.ExtensionsByType("text/html")
	mime.AddExtensionType(".custom", "application/x-custom")
	_ = mime.FormatMediaType("text/html", map[string]string{"charset": "utf-8"})
	_, _, _ = mime.ParseMediaType("text/html; charset=utf-8")
}

func main() {
	mimeExercise()
	fmt.Println("mime e2e test done")
}
