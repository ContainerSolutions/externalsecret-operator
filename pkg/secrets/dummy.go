package secrets

// DummySecretsBackend is a fake secrets backend for testing purposes
type DummySecretsBackend struct {
	suffix string
}

func init() {
	BackendRegister("dummy", NewDummySecretsBackend)
}

// NewDummySecretsBackend gives you an new DummySecretsBackend
func NewDummySecretsBackend() Backend {
	return &DummySecretsBackend{}
}

// Init implements SecretsBackend interface, sets the suffix
func (d *DummySecretsBackend) Init(parameters map[string]string) error {
	d.suffix = parameters["suffix"]
	return nil
}

// Get a key and returns a fake secrets key + suffix
func (d *DummySecretsBackend) Get(key string) (string, error) {
	return key + d.suffix, nil
}
