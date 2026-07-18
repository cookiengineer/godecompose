package main

import (
	"fmt"
	"regexp"
)

func regexpExercise() {
	re := regexp.MustCompile(`\w+`)
	_ = re.MatchString("hello")
	_ = re.FindString("hello world")
	_ = re.FindAllString("hello world", -1)
	_ = re.ReplaceAllString("hello world", "x")
	_ = re.Split("hello,world", -1)
}

func main() {
	regexpExercise()
	fmt.Println("regexp e2e test done")
}
