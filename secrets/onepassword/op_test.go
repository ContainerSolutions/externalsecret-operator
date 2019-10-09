package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
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

func TestErrItemInvalid(t *testing.T) {
	err := &ErrItemInvalid{item: "myitem"}

	expected := "1Password item 'myitem' is invalid. it should have a section 'External Secret Operator' with a field equal to the name of the item, 'myitem', and a value equal to the secret"

	actual := err.Error()
	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got '%s'", expected, actual)
	}
}

func TestAuthenticate_Err(t *testing.T) {
	op := &Op{GetterBuilder: &MockGetterBuilder{}}

	err := op.Authenticate("domain", "email", "secretKey", "masterPassword")
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
	expected := value

	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got '%s'", expected, actual)
	}
}

func TestGetItem_ErrItemInvalid_FailedGetItemMap(t *testing.T) {
	op := &Op{Getter: &MockGetter{}}

	_, err := op.GetItem("vault", "item")
	if err == nil {
		t.Fail()
	}
}

func TestGetItem_ErrItemInvald_MissingSection(t *testing.T) {
	itemMap := make(op.ItemMap)

	op := &Op{GetterBuilder: &MockGetterBuilder{itemMap: itemMap}}

	_ = op.Authenticate("domain", "email", "secretKey", "masterPassword")
	_, err := op.GetItem("vault", "item")

	if err == nil {
		t.Fail()
	}
}

func TestGetItem_ErrItemInvalid_MissingField(t *testing.T) {
	itemMap := make(op.ItemMap)
	fm := make(op.FieldMap)
	itemMap[op.SectionName(requiredSection)] = fm

	op := &Op{Getter: &MockGetter{itemMap: itemMap}}

	_, err := op.GetItem("vault", "item")

	switch err.(type) {
	case *ErrItemInvalid:
	default:
		t.Fail()
	}
}

func TestNotAuthenticatedGetItemMap(t *testing.T) {
	notAuthGetter := &NotAuthenticatedGetter{}

	_, err := notAuthGetter.GetItemMap(op.VaultName("vault"), op.ItemName("item"))
	actual := err.Error()
	expected := "failed to get an item map because you are not authenticated"
	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got '%s'", expected, actual)
	}
}
