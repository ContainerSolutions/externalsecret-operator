package dummy

import (
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
)

// Backend is a fake secrets backend for testing purposes
type Backend struct {
	suffix string
}

func init() {
	secrets.Register("dummy", New)
}

// New gives you an new Dummy Backend
func New() secrets.Backend {
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
