// Package onepassword implements a secrets backend for One Password.
package onepassword

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ContainerSolutions/externalsecret-operator/secrets/backend"
	"github.com/tidwall/gjson"
)

// Backend represents a Backend for onepassword
type Backend struct {
	Client OnePasswordClient
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
	backend.Client = OnePasswordCliClient{}
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

	session, err := b.Client.SignIn(parameters["domain"], parameters["email"], parameters["secretKey"], parameters["masterPassword"])
	if err != nil {
		return err
	} else {
		os.Setenv(session.Key, session.Value)
	}

	return nil
}

// Get retrieves the 1password item whose name matches the key and return the
// value of the 'password' field.
func (b *Backend) Get(key string) (string, error) {
	fmt.Println("Retrieving 1password item '" + key + "'.")

	item := b.Client.Get(key)
	if item == "" {
		return "", fmt.Errorf("Could not retrieve 1password item '" + key + "'.")
	}

	value := gjson.Get(item, "details.fields.#[name==\"password\"].value")
	if !value.Exists() {
		return "", fmt.Errorf("1password item '" + key + "' does not have a 'password' field.")
	}

	fmt.Println("1password item '" + key + "' value of 'password' field retrieved successfully.")

	return value.String(), nil
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
