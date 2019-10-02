// Package onepassword implements a secrets backend for One Password.
package onepassword

import (
	"fmt"

	"github.com/containersolutions/externalsecret-operator/secrets/backend"
	"github.com/pkg/errors"
)

type ErrSigninFailed struct {
	message string
}

func NewErrSigninFailed(message string) *ErrSigninFailed {
	return &ErrSigninFailed{
		message: message,
	}
}

func (e *ErrSigninFailed) Error() string {
	return "could not sign in to 1password: " + e.message
}

type ErrParameterMissing struct {
	parameter string
}

func (e *ErrParameterMissing) Error() string {
	return fmt.Sprintf("error reading 1password backend parameters: invalid init parameters: expected `%s` not found", e.parameter)
}

func NewErrParameterMissing(parameter string) *ErrParameterMissing {
	return &ErrParameterMissing{
		parameter: parameter,
	}
}

type ErrGetItem struct {
	itemName string
}

func (e *ErrGetItem) Error() string {
	return fmt.Sprintf("error retrieving 1password item '%s'", e.itemName)
}

func NewErrGetItem(itemName string) error {
	return &ErrGetItem{
		itemName: itemName,
	}
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
	sectionName         = "External Secret Operator"
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
	backend.OnePassword = &Cli{Op: &RealOp{}}
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

	err = b.OnePassword.SignIn(parameters[paramDomain], parameters[paramEmail], parameters[paramSecretKey], parameters[paramMasterPassword])
	if err != nil {
		return NewErrSigninFailed(err.Error())
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
		return "", NewErrGetItem(key)
	}

	return item, nil
}

func validateParameters(parameters map[string]string) error {
	for _, key := range paramKeys {
		_, found := parameters[key]
		if !found {
			return NewErrParameterMissing(key)
		}
	}

	return nil
}

type OnePassword interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) error
	GetItem(vault string, item string) (string, error)
}
