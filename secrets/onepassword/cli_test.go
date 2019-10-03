package onepassword

import (
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

func TestGetItem_ErrItemInvald(t *testing.T) {
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
