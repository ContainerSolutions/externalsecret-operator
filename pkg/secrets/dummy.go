package secrets

// DummySecretsBackend is a fake secrets backend for testing purposes
type DummySecretsBackend struct {
	Backend
	suffix string
}

func init() {
	BackendRegister("dummy", NewDummySecretsBackend)
}

// NewDummySecretsBackend gives you an new DummySecretsBackend
func NewDummySecretsBackend() BackendIface {
	return &DummySecretsBackend{}
}

// Init implements SecretsBackend interface, sets the suffix
func (d *DummySecretsBackend) Init(parameters ...interface{}) error {
	params := parameters[0].(map[string]string)
	d.suffix = params["suffix"]
	return nil
}

// Get a key and returns a fake secrets key + suffix
func (d *DummySecretsBackend) Get(key string) (string, error) {
	return key + d.suffix, nil
}
