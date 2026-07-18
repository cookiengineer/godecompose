package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func pathExercise() {
	_ = filepath.Join("a", "b", "c")
	_ = filepath.Base("/etc/hosts")
	_ = filepath.Dir("/etc/hosts")
	_ = filepath.Ext("file.txt")
	_, _ = filepath.Abs(".")
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error { return nil })
	_, _ = filepath.Glob("*.go")
}

func main() {
	pathExercise()
	fmt.Println("path e2e test done")
}
