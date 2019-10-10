// Package onepassword implements a secrets backend for One Password.
package onepassword

import (
	"fmt"

	"github.com/containersolutions/externalsecret-operator/secrets/backend"
	"github.com/pkg/errors"
)

type ErrInitFailed struct {
	message string
}

func (e *ErrInitFailed) Error() string {
	return fmt.Sprintf("1password backend init failed: %s", e.message)
}

type ErrGet struct {
	itemName string
	message  string
}

func (e *ErrGet) Error() string {
	return fmt.Sprintf("1password backend get '%s' failed: %s", e.itemName, e.message)
}

var (
	backendName         = "onepassword"
	defaultVault        = "Personal"
	paramDomain         = "domain"
	paramEmail          = "email"
	paramSecretKey      = "secretKey"
	paramMasterPassword = "masterPassword"
	paramVault          = "vault"
	paramKeys           = []string{paramDomain, paramEmail, paramSecretKey, paramMasterPassword, paramVault}
	errSigninFailed     = errors.New("could not sign in to 1password")
)

// Backend implementation for 1Password
type Backend struct {
	OnePassword OnePassword
	Vault       string
}

func init() {
	backend.Register(backendName, NewBackend)
}

// NewBackend returns a 1Password backend
func NewBackend() backend.Backend {
	backend := &Backend{}
	backend.OnePassword = &Op{GetterBuilder: &OpGetterBuilder{}}
	backend.Vault = defaultVault
	return backend
}

// Init reads secrets from the parameters and sign in to 1password.
func (b *Backend) Init(parameters map[string]string) error {
	err := validateParameters(parameters)
	if err != nil {
		return err
	}
	b.Vault = parameters[paramVault]

	err = b.OnePassword.Authenticate(parameters[paramDomain], parameters[paramEmail], parameters[paramMasterPassword], parameters[paramSecretKey])
	if err != nil {
		return &ErrInitFailed{message: err.Error()}
	}
	fmt.Println("signed into 1password successfully")

	return nil
}

// Get retrieves the 1password item whose name matches the key and return the
// value of the 'password' field.
func (b *Backend) Get(key string) (string, error) {
	fmt.Println("Retrieving 1password item '" + key + "'.")

	item, err := b.OnePassword.GetItem(b.Vault, key)
	if err != nil {
		return "", &ErrGet{itemName: key, message: err.Error()}
	}

	return item, nil
}

func validateParameters(parameters map[string]string) error {
	for _, key := range paramKeys {
		value, found := parameters[key]
		fmt.Printf("parameter '%s' has length: '%d'\n", key, len(value))

		if !found {
			return &ErrInitFailed{message: fmt.Sprintf("expected parameter '%s'", key)}
		} else if value == "" {
			return &ErrInitFailed{message: fmt.Sprintf("parameter '%s' is empty", key)}
		}
	}
	return nil
}

type OnePassword interface {
	Authenticate(domain string, email string, masterPassword string, secretKey string) error
	GetItem(vault string, item string) (string, error)
}
