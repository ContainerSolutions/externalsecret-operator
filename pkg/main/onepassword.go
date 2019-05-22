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
		fmt.Println(err.Error())
	}
}
