package onepassword

import (
	op "github.com/ameier38/onepassword"
)

type OnePassword interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) error
	GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error)
}

type OP struct {
	OP *op.Client
}

func (c *OP) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	op, err := op.NewClient("/usr/local/bin/op", domain, email, masterPassword, secretKey)
	c.OP = op
	return err
}

func (c *OP) GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error) {
	return c.OP.GetItem(vault, item)
}
