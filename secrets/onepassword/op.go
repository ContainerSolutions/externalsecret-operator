package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

type RealOp struct {
	client *op.Client
}

func (r *RealOp) NewClient(domain string, email string, masterPassword string, secretKey string) (*op.Client, error) {
	client, err := op.NewClient("/usr/local/bin/op", domain, email, masterPassword, secretKey)
	if err != nil {
		return nil, err
	}
	r.client = client
	return client, nil
}

func (r *RealOp) GetItem(vaultName op.VaultName, itemName op.ItemName) (op.ItemMap, error) {
	return r.client.GetItem(vaultName, itemName)
}

type FakeOp struct {
	VaultName string
	ItemName  string
	ItemValue string
	signInOk  bool
}

func (f *FakeOp) GetItem(vaultName op.VaultName, itemName op.ItemName) (op.ItemMap, error) {
	im := make(op.ItemMap)
	if string(itemName) == string(f.ItemName) {
		fm := make(op.FieldMap)
		fm[op.FieldName(f.ItemName)] = op.FieldValue(f.ItemValue)
		im[op.SectionName("External Secret Operator")] = fm
		return im, nil
	}
	return im, nil
}

func (f *FakeOp) NewClient(domain string, email string, masterPassword string, secretKey string) (*op.Client, error) {
	if !f.signInOk {
		return nil, fmt.Errorf("fake op sign in programmed to fail")
	}
	return nil, nil
}

func (f *FakeOp) SignIn(vaultName op.VaultName, itemName op.ItemName) error {
	if f.signInOk {
		return nil
	}
	return fmt.Errorf("fake op sign in programmed to fail")
}

func (f *FakeOp) SignInOk(signInOk bool) {
	f.signInOk = signInOk
}

func NewFakeOp(vaultName string, itemName string, itemValue string) *FakeOp {
	return &FakeOp{VaultName: vaultName, ItemName: itemName, ItemValue: itemValue}
}
