package onepassword

import (
	op "github.com/ameier38/onepassword"
)

type Cli interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) error
	GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error)
}

type OPCli struct {
	OP *op.Client
}

func (c *OPCli) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	op, err := op.NewClient("/usr/local/bin/op", domain, email, masterPassword, secretKey)
	c.OP = op
	return err
}

func (c *OPCli) GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error) {
	return c.OP.GetItem(vault, item)
}
