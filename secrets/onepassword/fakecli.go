package onepassword

import (
	op "github.com/ameier38/onepassword"
)

type FakeCli struct {
	ItemName  string
	ItemValue string
	SignInOK  bool
}

func (f *FakeCli) GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error) {
	im := make(op.ItemMap)

	fm := make(op.FieldMap)
	fm[op.FieldName(f.ItemName)] = op.FieldValue(f.ItemValue)

	im[op.SectionName("External Secret Operator")] = fm

	return im, nil
}

func (f *FakeCli) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	if f.SignInOK {
		return nil
	} else {
		return NewErrSigninFailed("fake cli configured for failure")
	}
}
