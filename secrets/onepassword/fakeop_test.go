package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
)

func TestFakeOPGetItem(t *testing.T) {
	vault := "vault"
	item := "item"
	value := "value"

	f := NewFakeOp(vault, item, value)

	vaultName := op.VaultName(vault)
	itemName := op.ItemName(item)

	itemMap, _ := f.GetItem(vaultName, itemName)

	sectionName := "External Secret Operator"
	if string(itemMap[op.SectionName(sectionName)][op.FieldName(itemName)]) != value {
		t.Fail()
		fmt.Printf("expected itemMap with item '%s' and value '%s' underneath section '%s' but got itemMap: '%v'", item, value, sectionName, itemMap)
	}
}

func TestFakeOPGetItem_ErrItemNotFound(t *testing.T) {
	vault := "vault"
	item := "item"
	value := "value"

	f := NewFakeOp(vault, item, value)

	vaultName := op.VaultName(vault)
	itemName := op.ItemName("nonExistentItem")

	itemMap, _ := f.GetItem(vaultName, "nonExistentItem")

	sectionName := "External Secret Operator"
	actual := itemMap[op.SectionName(sectionName)][op.FieldName(itemName)]
	if actual != "" {
		t.Fail()
		fmt.Printf("expected an empty string because item 'nonExistenItem' does not exist but got: '%s'", actual)
	}
}
