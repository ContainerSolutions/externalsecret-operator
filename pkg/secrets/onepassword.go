package secrets

import (
	"fmt"
)

type OnePasswordBackend struct {
	Backend
}

func NewOnePasswordBackend() *OnePasswordBackend {
	backend := &OnePasswordBackend{}
	backend.Init()
	return backend
}

func (b *OnePasswordBackend) Init(params ...interface{}) error {
	fmt.Println("Initializing 1password backend.")

	return nil
}

func (b *OnePasswordBackend) Get(key string) (string, error) {
	var value = "secretValue"

	return value, nil
}
