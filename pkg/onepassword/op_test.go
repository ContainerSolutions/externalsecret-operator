package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ContainerSolutions/onepassword"

	. "github.com/containersolutions/externalsecret-operator/pkg/internal"
)

type MockGetterBuilder struct {
	itemMap op.ItemMap
}

func (m *MockGetterBuilder) NewGetter(domain string, email string, masterPassword string, secretKey string) (Getter, error) {
	if m.itemMap == nil {
		return nil, fmt.Errorf("mock op: could not build new getter")
	} else {
		return &MockGetter{itemMap: m.itemMap}, nil
	}
}

type MockGetter struct {
	itemMap op.ItemMap
}

func (m *MockGetter) GetItemMap(op.VaultName, op.ItemName) (op.ItemMap, error) {
	if m.itemMap == nil {
		return nil, fmt.Errorf("mock op: could not get item")
	}
	return m.itemMap, nil
}

func TestErrMissingSection(t *testing.T) {
	err := &ErrMissingSection{item: "myitem"}

	expected := "missing section 'External Secret Operator' in 1Password item 'myitem'"

	AssertEquals(t, expected, err.Error())
}

func TestAuthenticate(t *testing.T) {
	itemMap := make(op.ItemMap)

	builder := &MockGetterBuilder{itemMap: itemMap}
	op := &Op{GetterBuilder: builder}

	op.Authenticate("domain", "email", "masterPassword", "secretKey")

	AssertNotNil(t, op.Getter)
}

func TestAuthenticate_Err(t *testing.T) {
	op := &Op{GetterBuilder: &MockGetterBuilder{}}

	err := op.Authenticate("domain", "email", "masterPassword", "secretKey")
	if err == nil {
		t.Fail()
	}
}

func TestGetItem(t *testing.T) {
	item := "item"
	value := "value"
	vault := "vault"

	itemMap := make(op.ItemMap)
	fm := make(op.FieldMap)
	fieldName := op.FieldName(item)
	fieldValue := op.FieldValue(value)

	fm[fieldName] = fieldValue
	itemMap[op.SectionName(requiredSection)] = fm

	op := &Op{Getter: &MockGetter{itemMap: itemMap}}

	actual, _ := op.GetItem(vault, item)

	AssertEquals(t, value, actual)
}

func TestGetItem_ErrFailedGetItemMap(t *testing.T) {
	op := &Op{Getter: &MockGetter{}}

	_, err := op.GetItem("vault", "item")

	expected := "failed to get itemMap of 1Password item 'item': mock op: could not get item"
	AssertEquals(t, expected, err.Error())
}

func TestGetItem_ErrMissingSection(t *testing.T) {
	item := "item"
	value := "value"

	itemMap := make(op.ItemMap)
	fm := make(op.FieldMap)
	fieldName := op.FieldName(item)
	fieldValue := op.FieldValue(value)

	fm[fieldName] = fieldValue

	op := &Op{Getter: &MockGetter{itemMap: itemMap}}

	_, err := op.GetItem("vault", "item")

	expected := "missing section 'External Secret Operator' in 1Password item 'item'"
	AssertEquals(t, expected, err.Error())
}

func TestGetItem_ErrMissingField(t *testing.T) {
	itemMap := make(op.ItemMap)
	fm := make(op.FieldMap)
	itemMap[op.SectionName(requiredSection)] = fm

	op := &Op{Getter: &MockGetter{itemMap: itemMap}}

	_, err := op.GetItem("vault", "item")

	expected := "missing field 'item' in 1Password item 'item'"
	AssertEquals(t, expected, err.Error())
}

func TestNotAuthenticatedGetItemMap(t *testing.T) {
	notAuthGetter := &NotAuthenticatedGetter{}

	_, err := notAuthGetter.GetItemMap(op.VaultName("vault"), op.ItemName("item"))

	expected := "failed to get an item map because you are not authenticated"
	AssertEquals(t, expected, err.Error())
}
