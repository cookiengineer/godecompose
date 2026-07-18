package main

import (
	"flag"
	"fmt"
)

func flagExercise() {
	name := flag.String("name", "default", "name flag")
	count := flag.Int("count", 0, "count flag")
	enabled := flag.Bool("enabled", false, "enabled flag")
	flag.Parse()
	_ = *name
	_ = *count
	_ = *enabled

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fsName := fs.String("name", "default", "name flag")
	fsCount := fs.Int("count", 0, "count flag")
	fsBool := fs.Bool("bool", false, "bool flag")
	_ = fsName
	_ = fsCount
	_ = fsBool
	_ = fs.Parse([]string{"-name=test"})
}

func main() {
	flagExercise()
	fmt.Println("flag e2e test done")
}
