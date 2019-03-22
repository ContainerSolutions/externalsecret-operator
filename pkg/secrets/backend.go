package secrets

// SecretsBackend is an interface to a Secret Store
type SecretsBackend interface {
	Init(...interface{}) error
	Get(string) (string, error)
}
