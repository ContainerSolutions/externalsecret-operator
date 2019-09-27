package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
)

func TestGetItem(t *testing.T) {
	vault := "Shared"
	secretKey := "secretKey"
	secretValue := "secretValue"
	sectionName := "External Secret Operator"

	cli := &FakeCli{
		Key:   secretKey,
		Value: secretValue,
	}

	itemMap, _ := cli.GetItem(op.VaultName(vault), op.ItemName(secretKey))

	actual := itemMap[op.SectionName(sectionName)][op.FieldName(secretKey)]

	if actual != op.FieldValue(secretValue) {
		t.Fail()
		fmt.Printf("Expected to retrieve secretValue. Got '%s' wanted '%s'", actual, secretValue)
	}
}
