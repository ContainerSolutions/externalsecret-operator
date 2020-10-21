package utils

import (
	"crypto/rand"
	"math/big"
)

const validObjChars = "0123456789abcdefghijklmnopqrstuvwxyz"

// RandomBytes generate random bytes
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandomInt returns a random int64
func RandomInt() (int64, error) {
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(validObjChars))))
	if err != nil {
		return 0, err
	}

	return randomInt.Int64(), nil
}

// RandomStringObjectSafe returns a random string that is safe to use as and k8s object Name
//  https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
func RandomStringObjectSafe(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}

	for i := range b {
		randomInt, err := RandomInt()
		if err != nil {
			return "", err
		}
		b[i] = validObjChars[randomInt]
	}
	return string(b), nil

}
