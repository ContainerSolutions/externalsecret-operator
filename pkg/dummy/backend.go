// Package dummy implements an example backend that can be used for testing
// purposes. It acceps a "suffix" as a parameter that will be appended to the
// key passed to the Get function.
package dummy

import (
	"fmt"

	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("dummy")

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
	if len(parameters) == 0 {
		log.Error(fmt.Errorf("error"), "empty or invalid parameters: ")
		return fmt.Errorf("empty or invalid parameters")
	}

	suffix, ok := parameters["Suffix"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing parameters: ")
		return fmt.Errorf("missing parameters")
	}

	d.suffix = suffix
	return nil
}

// Get a key and returns a fake secrets key + suffix
func (d *Backend) Get(key string, version string) (string, error) {
	if d.suffix == "" {
		return "", fmt.Errorf("backend is not initialized")
	}
	return key + version + d.suffix, nil
}
