package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func rsaExercise() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	_ = key

	msg := []byte("hello")
	hash := sha256.Sum256(msg)

	enc, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, msg, nil)
	_, _ = rsa.DecryptOAEP(sha256.New(), rand.Reader, key, enc, nil)

	sig, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	_ = sig
	_ = rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hash[:], sig)
}

func main() {
	rsaExercise()
	fmt.Println("rsa e2e test done")
}
