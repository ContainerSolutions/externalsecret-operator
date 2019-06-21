package main

import (
	"fmt"
	"os"

	"github.com/ContainerSolutions/externalsecret-operator/secrets/onepassword"
)

func main() {
	domain := os.Getenv("OP_DOMAIN")
	email := os.Getenv("OP_EMAIL")
	secretKey := os.Getenv("OP_SECRET_KEY")
	masterPassword := os.Getenv("OP_MASTER_PASSWORD")

	op := onepassword.OP{}

	stdout, err := op.SignIn(domain, email, secretKey, masterPassword)
	if err != nil {
		fmt.Println("could not signin via 'op'", err)
	}

	fmt.Println(stdout)
}
