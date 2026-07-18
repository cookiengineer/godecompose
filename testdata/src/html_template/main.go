package main

import (
	"fmt"
	"html/template"
	"os"
)

func htmlTemplateExercise() {
	t := template.New("test")
	t, _ = t.Parse("<h1>{{.Title}}</h1>")
	_ = t.Execute(os.Stdout, map[string]string{"Title": "hello"})
	_ = t.ExecuteTemplate(os.Stdout, "test", nil)

	_ = template.Must(t, nil)

	t2, _ := template.ParseFiles("/nonexistent")
	_ = t2

	t3, _ := template.ParseGlob("/nonexistent/*.html")
	_ = t3
}

func main() {
	htmlTemplateExercise()
	fmt.Println("html/template e2e test done")
}
