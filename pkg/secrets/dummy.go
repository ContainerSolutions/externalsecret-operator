package secrets

type ContextKey string

// DummySecretsBackend is a fake secrets backend for testing purposes
type DummySecretsBackend struct {
	suffix string
}

// NewDummySecretsBackend gives you an new DummySecretsBackend
func NewDummySecretsBackend() *DummySecretsBackend {
	return &DummySecretsBackend{}
}

// Init implements SecretsBackend interface, sets the suffix
func (d *DummySecretsBackend) Init(parameters ...interface{}) error {
	d.suffix = parameters[0].(string)
	return nil
}

// Get a key and returns a fake secrets key + suffix
func (d *DummySecretsBackend) Get(key string) (string, error) {
	return key + d.suffix, nil
}
