package onepassword

import (
	"fmt"
	"testing"

	op "github.com/ameier38/onepassword"
)

func TestGetItem(t *testing.T) {
	vault := "Shared"
	itemName := "itemName"
	itemValue := "itemValue"
	sectionName := "External Secret Operator"

	onepassword := &FakeOnePassword{
		ItemName:  itemName,
		ItemValue: itemValue,
		SignInOK:  true,
	}

	itemMap, _ := onepassword.GetItem(op.VaultName(vault), op.ItemName(itemName))

	actual := itemMap[op.SectionName(sectionName)][op.FieldName(itemName)]

	if actual != op.FieldValue(itemValue) {
		t.Fail()
		fmt.Printf("Expected to retrieve item value. Got '%s' wanted '%s'", actual, itemValue)
	}
}

func TestGetItem_ErrGetItem(t *testing.T) {
	vault := "Shared"
	itemName := "itemName"
	itemValue := "itemValue"

	onepassword := &FakeOnePassword{
		ItemName:  itemName,
		ItemValue: itemValue,
		SignInOK:  true,
	}

	nonExistentItem := "nonExistentItem"

	_, err := onepassword.GetItem(op.VaultName(vault), op.ItemName(nonExistentItem))
	if err == nil {
		t.Fail()
		fmt.Printf("expected error getting item")
	}
}

func TestSignIn(t *testing.T) {
	itemName := "itemName"
	itemValue := "itemValue"

	onepassword := &FakeOnePassword{
		ItemName:  itemName,
		ItemValue: itemValue,
		SignInOK:  true,
	}

	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"

	err := onepassword.SignIn(domain, email, secretKey, masterPassword)
	if err != nil {
		t.Fail()
		fmt.Println("signin should be successful")
	}
}
