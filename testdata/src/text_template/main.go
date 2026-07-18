package main

import (
	"fmt"
	"os"
	"text/template"
)

func textTemplateExercise() {
	t := template.New("test")
	t, _ = t.Parse("{{.Name}}")
	_ = t.Execute(os.Stdout, map[string]string{"Name": "world"})
	_ = t.ExecuteTemplate(os.Stdout, "test", nil)

	_ = template.Must(t, nil)
}

func main() {
	textTemplateExercise()
	fmt.Println("text/template e2e test done")
}
