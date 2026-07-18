package main

import (
	"errors"
	"fmt"
)

var ErrTest = errors.New("test error")

func errorsExercise() {
	err := errors.New("something failed")
	_ = err

	if errors.Is(err, ErrTest) {
		fmt.Println("matched")
	}

	var pathErr *error
	_ = pathErr
	_ = errors.As(err, &pathErr)

	unwrapped := errors.Unwrap(err)
	_ = unwrapped

	_ = errors.Join(err, ErrTest)
}

func main() {
	errorsExercise()
	fmt.Println("errors e2e test done")
}
