package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

type ErrItemInvalid struct {
	item    string
	section string
}

func (e *ErrItemInvalid) Error() string {
	return fmt.Sprintf(
		`item '%s' is invalid.
		 a 1Password item should have a section called '%s'
		 with a field equal to the name of the item
		 and a value equal to the secret you want to store`, e.item, e.section)
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

	section := "External Secret Operator"
	sectionMap := itemMap[op.SectionName(section)]

	if sectionMap == nil {
		return "", &ErrItemInvalid{item: item, section: section}
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
