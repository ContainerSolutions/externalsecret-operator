package main

import (
	"fmt"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
)

func main() {
	vault := "Personal"
	client := secrets.OnePasswordCliClient{}
	backend := secrets.NewOnePasswordBackend(vault, client)

	err := backend.Init(vault)
	if err != nil {
		fmt.Println("Init: " + err.Error())
	}

	key := "testkey"
	value, err := backend.Get(key)
	if err != nil {
		fmt.Println("Get: " + err.Error())
	}

	fmt.Println("Get '" + key + "' from vault '" + vault + "' = '" + value + "'")
}
