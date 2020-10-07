package utils

import (
	"math/rand"
)

// RandomString generate random lowercase-alphanumeric subdomain valid value
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
