package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

var requiredSection = "External Secret Operator"

type ErrItemInvalid struct {
	item string
}

func (e *ErrItemInvalid) Error() string {
	return fmt.Sprintf("1Password item '%s' is invalid. section '%s' is missing. this section should contain a field equal to the name of the item, '%s', and a value equal to the secret", e.item, requiredSection, e.item)
}

type Cli struct {
	Op Op
}

func (c Cli) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	_, err := c.Op.NewClient(domain, email, masterPassword, secretKey)
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

	sectionMap := itemMap[op.SectionName(requiredSection)]
	if sectionMap == nil {
		return "", &ErrItemInvalid{item: item}
	}

	itemValue := sectionMap[op.FieldName(op.ItemName(item))]
	if itemValue == "" {
		return "", fmt.Errorf("expected section '%s' to have an field with name '%s'", sectionName, item)
	}

	return string(itemValue), nil
}

type Op interface {
	NewClient(domain string, email string, masterPassword string, secretKey string) (*op.Client, error)
	GetItem(op.VaultName, op.ItemName) (op.ItemMap, error)
}
