package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

var executablePath = "/usr/local/bin/op"
var requiredSection = "External Secret Operator"

type ErrItemInvalid struct {
	item string
}

func (e *ErrItemInvalid) Error() string {
	return fmt.Sprintf("1Password item '%s' is invalid. it should have a section '%s' with a field equal to the name of the item, '%s', and a value equal to the secret", e.item, requiredSection, e.item)
}

type Getter interface {
	GetItemMap(vaultName op.VaultName, itemName op.ItemName) (op.ItemMap, error)
}

type OpGetter struct {
	client *op.Client
}

func (o OpGetter) GetItemMap(vaultName op.VaultName, itemName op.ItemName) (op.ItemMap, error) {
	return o.client.GetItem(vaultName, itemName)
}

type NotAuthenticatedGetter struct{}

func (n NotAuthenticatedGetter) GetItemMap(vault op.VaultName, itemName op.ItemName) (op.ItemMap, error) {
	return nil, fmt.Errorf("failed to get an item map because you are not authenticated")
}

type GetterBuilder interface {
	NewGetter(domain, email, masterPassword, secretKey string) (Getter, error)
}

type OpGetterBuilder struct{}

func (o OpGetterBuilder) NewGetter(domain, email, masterPassword, secretKey string) (Getter, error) {
	client, err := op.NewClient(executablePath, domain, email, masterPassword, secretKey)
	if err != nil {
		return &NotAuthenticatedGetter{}, err
	}
	return &OpGetter{client: client}, nil
}

type Op struct {
	Getter        Getter
	GetterBuilder GetterBuilder
}

func (o *Op) Authenticate(domain string, email string, secretKey string, masterPassword string) error {
	getter, err := o.GetterBuilder.NewGetter(domain, email, secretKey, masterPassword)
	if err != nil {
		return err
	}
	o.Getter = getter
	return nil
}

func (o *Op) GetItem(vault string, item string) (string, error) {
	itemMap, err := o.Getter.GetItemMap(op.VaultName(vault), op.ItemName(item))
	if err != nil {
		return "", err
	}

	sectionMap := itemMap[op.SectionName(requiredSection)]
	if sectionMap == nil {
		return "", &ErrItemInvalid{item: item}
	}

	itemValue := sectionMap[op.FieldName(op.ItemName(item))]
	if itemValue == "" {
		return "", &ErrItemInvalid{item: item}
	}

	return string(itemValue), nil
}
