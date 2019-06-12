// Package dummy implements an example backend that can be used for testing
// purposes. It acceps a "suffix" as a parameter that will be appended to the
// key passed to the Get function.
package dummy

import (
	"github.com/ContainerSolutions/externalsecret-operator/secrets"
)

// Backend is a fake secrets backend for testing purposes
type Backend struct {
	suffix string
}

func init() {
	secrets.Register("dummy", NewBackend)
}

// NewBackend gives you an NewBackend Dummy Backend
func NewBackend() secrets.Backend {
	return &Backend{}
}

// Init implements SecretsBackend interface, sets the suffix
func (d *Backend) Init(parameters map[string]string) error {
	d.suffix = parameters["suffix"]
	return nil
}

// Get a key and returns a fake secrets key + suffix
func (d *Backend) Get(key string) (string, error) {
	return key + d.suffix, nil
}
