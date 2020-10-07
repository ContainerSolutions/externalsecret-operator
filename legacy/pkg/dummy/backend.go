// Package dummy implements an example backend that can be used for testing
// purposes. It acceps a "suffix" as a parameter that will be appended to the
// key passed to the Get function.
package dummy

import "github.com/containersolutions/externalsecret-operator/pkg/backend"

// Backend is a fake secrets backend for testing purposes
type Backend struct {
	suffix string
}

func init() {
	backend.Register("dummy", NewBackend)
}

// NewBackend gives you an NewBackend Dummy Backend
func NewBackend() backend.Backend {
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
