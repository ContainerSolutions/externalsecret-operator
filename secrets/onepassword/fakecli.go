package onepassword

import (
	op "github.com/ameier38/onepassword"
)

type FakeCli struct {
	Key   string
	Value string
}

func (f *FakeCli) GetItem(vault op.VaultName, item op.ItemName) (op.ItemMap, error) {
	im := make(op.ItemMap)

	fm := make(op.FieldMap)
	fm[op.FieldName(f.Key)] = op.FieldValue(f.Value)

	im[op.SectionName("External Secret Operator")] = fm	

	return im, nil
}
