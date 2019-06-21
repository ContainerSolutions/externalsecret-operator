package main

import (
	"fmt"
	"os"

	op "github.com/ameier38/onepassword"
)

func main() {
	// domain := os.Getenv("OP_DOMAIN")
	email := os.Getenv("OP_EMAIL")
	secretKey := os.Getenv("OP_SECRET_KEY")
	masterPassword := os.Getenv("OP_MASTER_PASSWORD")

	client, err := op.NewClient("/usr/local/bin/op", "containersolutions", email, masterPassword, secretKey)
	if err != nil {
		fmt.Println("could not signin via 'op'", err)
	}

	itemMap, err := client.GetItem(op.VaultName("test vault one"), op.ItemName("testkey"))
	if itemMap == nil {
		fmt.Println(fmt.Errorf("could not retrieve 1password item 'testkey'."))
	}
	if err != nil {
		fmt.Println(fmt.Errorf("error retrieving 1password item 'testkey'."))
	}

	fmt.Println(itemMap)
	fmt.Println(itemMap[""])
	fmt.Println(itemMap[""]["testkey"])
	fmt.Println(itemMap["testkey"])
}
