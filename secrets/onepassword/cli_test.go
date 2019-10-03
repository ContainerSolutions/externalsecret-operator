package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
)

type MockOp struct {
	itemMap op.ItemMap
}

func (m *MockOp) NewClient(domain string, email string, secretKey string, masterPassword string) (*op.Client, error) {
	return nil, &ErrOpNewClient{message: "op: could not create new client"}
}

func (m *MockOp) GetItem(op.VaultName, op.ItemName) (op.ItemMap, error) {
	if m.itemMap == nil {
		return nil, &ErrOpGetItem{message: "op: could not get item"}
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

func TestSignIn_Err(t *testing.T) {
	cli := &Cli{Op: &MockOp{}}

	err := cli.SignIn("domain", "email", "secretKey", "masterPassword")

	switch err.(type) {
	case *ErrOpNewClient:
	default:
		t.Fail()
	}
}

func TestGetItem_ErrGetItem(t *testing.T) {
	cli := &Cli{Op: &MockOp{}}

	_, err := cli.GetItem("vault", "item")

	switch err.(type) {
	case *ErrOpGetItem:
	default:
		t.Fail()
	}
}

func TestGetItem_ErrItemInvald_MissingSection(t *testing.T) {
	mockOp := &MockOp{}
	mockOp.itemMap = make(op.ItemMap)
	cli := &Cli{Op: mockOp}

	_, err := cli.GetItem("vault", "item")

	switch err.(type) {
	case *ErrItemInvalid:
	default:
		t.Fail()
	}
}

func TestGetItem_ErrItemInvalid_MissingField(t *testing.T) {
	mockOp := &MockOp{}
	itemMap := make(op.ItemMap)
	fm := make(op.FieldMap)
	itemMap[op.SectionName(requiredSection)] = fm
	mockOp.itemMap = itemMap
	cli := &Cli{Op: mockOp}

	_, err := cli.GetItem("vault", "item")

	switch err.(type) {
	case *ErrItemInvalid:
	default:
		t.Fail()
	}
}
