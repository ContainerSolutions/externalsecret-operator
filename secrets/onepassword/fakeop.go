package onepassword

import op "github.com/ameier38/onepassword"

type FakeOp struct {
	VaultName string
	ItemName  string
	ItemValue string
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

func NewFakeOp(vaultName string, itemName string, itemValue string) *FakeOp {
	return &FakeOp{VaultName: vaultName, ItemName: itemName, ItemValue: itemValue}
}
