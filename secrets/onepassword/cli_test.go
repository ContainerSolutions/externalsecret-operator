package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
)

type MockOp struct{}

func (m *MockOp) NewClient(domain string, email string, secretKey string, masterPassword string) (*op.Client, error) {
	return nil, fmt.Errorf("could not create new `op` client")
}

func (m *MockOp) GetItem(op.VaultName, op.ItemName) (op.ItemMap, error) {
	return nil, nil
}

func TestSignIn_Err(t *testing.T) {
	cli := &Cli{Op: &MockOp{}}

	err := cli.SignIn("domain", "email", "secretKey", "masterPassword")

	expected := "could not create new `op` client"
	actual := err.Error()
	if actual != expected {
		t.Fail()
		fmt.Printf("expected '%s' got '%s'", expected, actual)
	}
}
