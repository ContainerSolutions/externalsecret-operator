package onepassword

import (
	"fmt"

	op "github.com/ameier38/onepassword"
)

type FakeOnePassword struct {
	ItemName  string
	ItemValue string
	SignInOK  bool
}

func (f *FakeOnePassword) GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error) {
	if string(item) == f.ItemName {
		im := make(op.ItemMap)

		fm := make(op.FieldMap)
		fm[op.FieldName(f.ItemName)] = op.FieldValue(f.ItemValue)

		im[op.SectionName("External Secret Operator")] = fm

		return im, nil
	} else {
		return nil, fmt.Errorf("could not retrieve item '%s'", string(item))
	}
}

func (f *FakeOnePassword) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	if f.SignInOK {
		return nil
	} else {
		return NewErrSigninFailed("fake cli configured for failure")
	}
}
