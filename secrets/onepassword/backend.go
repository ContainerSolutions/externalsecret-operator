// Package onepassword implements a secrets backend for One Password.
package onepassword

import (
	"fmt"
	"reflect"

	"github.com/ContainerSolutions/externalsecret-operator/secrets/backend"
	op "github.com/ameier38/onepassword"
)

// Backend represents a Backend for onepassword
type Backend struct {
	Client *op.Client
	Vault  string
}

type Session struct {
	Key   string
	Value string
}

func init() {
	backend.Register("onepassword", NewBackend)
}

// NewBackend returns a Backend for onepassword
func NewBackend() backend.Backend {
	backend := &Backend{}
	backend.Vault = "Personal"
	return backend
}

// Init reads secrets from the parameters and sign in to 1password.
func (b *Backend) Init(parameters map[string]string) error {

	err := validateParameters(parameters)
	if err != nil {
		return fmt.Errorf("Error reading 1password backend parameters: %v", err)
	}

	b.Vault = parameters["vault"]

	client, err := op.NewClient("/usr/local/bin/op", parameters["domain"], parameters["email"], parameters["masterPassword"], parameters["secretKey"])
	if err != nil {
		fmt.Println(fmt.Sprintf("could not sign in to 1password %s", err))
	}
	b.Client = client

	return nil
}

// Get retrieves the 1password item whose name matches the key and return the
// value of the 'password' field.
func (b *Backend) Get(key string) (string, error) {
	fmt.Println("Retrieving 1password item '" + key + "'.")

	itemMap, err := b.Client.GetItem(op.VaultName(b.Vault), op.ItemName(key))
	if itemMap != nil {
		return "", fmt.Errorf("could not retrieve 1password item '" + key + "'.")
	}
	if err != nil {
		return "", fmt.Errorf("error retrieving 1password item '" + key + "'.")
	}

	return string(itemMap["externalsecretoperator"]["testkey"]), nil
}

func validateParameters(parameters map[string]string) error {

	paramKeys := []string{"domain", "email", "secretKey", "masterPassword", "vault"}

	for _, key := range paramKeys {
		paramValue, found := parameters[key]
		if !found {
			return fmt.Errorf("invalid init parameters: expected `%v` not found", key)
		}

		paramType := reflect.TypeOf(paramValue)
		if paramType.Kind() != reflect.String {
			return fmt.Errorf("invalid init parameters: expected `%v` of type `string` got `%v`", key, paramType)
		}
	}

	return nil
}
