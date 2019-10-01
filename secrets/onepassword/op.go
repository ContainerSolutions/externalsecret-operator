package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

type FakeOp struct {
	VaultName string
	ItemName  string
	ItemValue string
	SignInOk  bool
}

func (f FakeOp) GetItem(vaultName op.VaultName, itemName op.ItemName) (op.ItemMap, error) {
	im := make(op.ItemMap)
	if string(itemName) == string(f.ItemName) {
		fm := make(op.FieldMap)
		fm[op.FieldName(f.ItemName)] = op.FieldValue(f.ItemValue)
		im[op.SectionName("External Secret Operator")] = fm
		return im, nil
	}
	return im, nil
}

func (f FakeOp) NewClient(executable string, domain string, email string, masterPassword string, secretKey string) (*op.Client, error) {
	if !f.SignInOk {
		return nil, fmt.Errorf("fake op sign in programmed to fail")
	}
	return nil, nil
}

func (f FakeOp) SignIn(vaultName op.VaultName, itemName op.ItemName) error {
	if f.SignInOk {
		return nil
	}
	return fmt.Errorf("fake op sign in programmed to fail")
}

func NewFakeOp(vaultName string, itemName string, itemValue string) *FakeOp {
	return &FakeOp{VaultName: vaultName, ItemName: itemName, ItemValue: itemValue}
}
