package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

var executablePath = "/usr/local/bin/op"
var requiredSection = "External Secret Operator"

type ErrFailedGetItemMap struct {
	message string
	item    string
}

func (e *ErrFailedGetItemMap) Error() string {
	return fmt.Sprintf("failed to get itemMap of 1Password item '%s': %s", e.item, e.message)
}

type ErrMissingField struct {
	item  string
	field string
}

func (e *ErrMissingField) Error() string {
	return fmt.Sprintf("missing field '%s' in 1Password item '%s'", e.field, e.item)
}

type ErrMissingSection struct {
	item string
}

func (e *ErrMissingSection) Error() string {
	return fmt.Sprintf("missing section '%s' in 1Password item '%s'", requiredSection, e.item)
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

func (o *Op) Authenticate(domain, email, masterPassword, secretKey string) error {
	getter, err := o.GetterBuilder.NewGetter(domain, email, masterPassword, secretKey)
	if err != nil {
		return err
	}
	o.Getter = getter
	return nil
}

func (o *Op) GetItem(vault string, item string) (string, error) {
	itemMap, err := o.Getter.GetItemMap(op.VaultName(vault), op.ItemName(item))
	if err != nil {
		return "", &ErrFailedGetItemMap{item: item, message: err.Error()}
	}

	sectionMap := itemMap[op.SectionName(requiredSection)]
	if sectionMap == nil {
		return "", &ErrMissingSection{item: item}
	}

	itemValue := sectionMap[op.FieldName(op.ItemName(item))]
	if itemValue == "" {
		return "", &ErrMissingField{item: item, field: item}
	}

	return string(itemValue), nil
}
