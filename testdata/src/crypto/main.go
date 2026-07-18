package main

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func cryptoExercise() {
	_ = sha256.Sum256([]byte("hello"))
	_ = sha256.New()
	_ = md5.Sum([]byte("hello"))
	_, _ = aes.NewCipher(make([]byte, 16))
	_ = hmac.New(sha256.New, make([]byte, 16))

	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
}

func main() {
	cryptoExercise()
	fmt.Println("crypto e2e test done")
}
