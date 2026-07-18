package main

import (
	"log"
	"os"
)

func logExercise() {
	log.Println("hello world")
	log.Printf("value: %d", 42)
	log.Fatal("fatal message")
	log.Fatalf("fatal: %s", "reason")
	log.Panic("panic message")
	log.Panicf("panic: %s", "reason")
}

func main() {
	os.Setenv("TEST_LOG", "0")
	logExercise()
}
