package main

import (
	"math/rand"
)

const shortCodealphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // the characters that can be used in the short code
func generateRandomshortCode(length int) string { // generate a random short code for the URL
	code := make([]byte, length) // created a byte slice of the specified length
	for i := range code {
		code[i] = shortCodealphabet[rand.Intn(len(shortCodealphabet))] // generate a random index in the shortCodealphabet and choose the character at that index to be the next character in the short code
	}
	return string(code)
}
