package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
	"github.com/pkg/errors"
)

// Client represents a 1Password Client
type Client interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) error
	Get(vault string, key string) (string, error)
}

// OP is a 1Password Client
type OP struct {
	OP *op.Client
}

// SignIn sings into a 1Password account corresponding to the parameters
func (client *OP) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	op, err := op.NewClient("/usr/local/bin/op", domain, email, masterPassword, secretKey)
	if err != nil {
		return errors.Wrap(err, "op signin failed")
	}

	client.OP = op

	return nil
}

// Get returns the value of the secret with 'key' or error
func (client *OP) Get(vault string, key string) (string, error) {
	itemMap, err := client.OP.GetItem(op.VaultName(vault), op.ItemName(key))
	if itemMap == nil {
		return "", fmt.Errorf("could not retrieve 1password item '%s'", key)
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("could not retrieve 1password item '%s'", key))
	}

	return string(itemMap["External Secret Operator"][op.FieldName("secretValue")]), nil
}
