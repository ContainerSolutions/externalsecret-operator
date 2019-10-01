package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

type Cli struct {
	Op Op
}

func (c Cli) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	_, err := c.Op.NewClient("/usr/local/bin/op", domain, email, masterPassword, secretKey)
	if err != nil {
		return err
	}
	return nil
}

func (c Cli) GetItem(vault string, item string) (string, error) {
	itemMap, err := c.Op.GetItem(op.VaultName(vault), op.ItemName(item))
	if err != nil {
		return "", err
	}

	sectionName := "External Secret Operator"
	sectionMap := itemMap[op.SectionName(sectionName)]

	if sectionMap == nil {
		return "", fmt.Errorf("expected item '%s' to have a section '%s'", item, sectionName)
	}

	itemValue := sectionMap[op.FieldName(op.ItemName(item))]

	if itemValue == "" {
		return "", fmt.Errorf("expected section '%s' to have an field with name '%s'", sectionName, item)
	}

	return string(itemValue), nil
}

type Op interface {
	NewClient(string, string, string, string, string) (*op.Client, error)
	GetItem(op.VaultName, op.ItemName) (op.ItemMap, error)
}
