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

	cli := &FakeCli{
		ItemName:  itemName,
		ItemValue: itemValue,
		SignInOK:  true,
	}

	itemMap, _ := cli.GetItem(op.VaultName(vault), op.ItemName(itemName))

	actual := itemMap[op.SectionName(sectionName)][op.FieldName(itemName)]

	if actual != op.FieldValue(itemValue) {
		t.Fail()
		fmt.Printf("Expected to retrieve item value. Got '%s' wanted '%s'", actual, itemValue)
	}
}

func TestSignIn(t *testing.T) {
	itemName := "itemName"
	itemValue := "itemValue"

	cli := &FakeCli{
		ItemName:  itemName,
		ItemValue: itemValue,
		SignInOK:  true,
	}

	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"

	err := cli.SignIn(domain, email, secretKey, masterPassword)
	if err != nil {
		t.Fail()
		fmt.Println("signin should be successful")
	}
}
