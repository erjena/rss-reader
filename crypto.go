package main

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createHash(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
