package secrets

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/kr/pty"
	"github.com/tidwall/gjson"
)

type OnePasswordBackend struct {
	OnePasswordClient
}

func NewOnePasswordBackend(vault string, client OnePasswordClient) *OnePasswordBackend {
	backend := &OnePasswordBackend{}
	backend.OnePasswordClient = client

	return backend
}

// Read secrets from the environment, sign in to 1password and clear the environment variables
func (b *OnePasswordBackend) Init(params ...interface{}) error {
	domain := os.Getenv("ONEPASSWORD_DOMAIN")
	if domain == "" {
		fmt.Println("Missing ONEPASSWORD_DOMAIN environment variable.")
		return fmt.Errorf("Missing ONEPASSWORD_DOMAIN environment variable.")
	}

	email := os.Getenv("ONEPASSWORD_EMAIL")
	if email == "" {
		fmt.Println("Missing ONEPASSWORD_EMAIL environment variable.")
		return fmt.Errorf("Missing ONEPASSWORD_EMAIL environment variable.")
	}

	secretKey := os.Getenv("ONEPASSWORD_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("Missing one or more ONEPASSWORD_SECRET_KEY environment variable.")
		return fmt.Errorf("Missing ONEPASSWORD_SECRET_KEY environment variable.")
	}

	masterPassword := os.Getenv("ONEPASSWORD_MASTER_PASSWORD")
	if masterPassword == "" {
		fmt.Println("Missing one or more ONEPASSWORD_MASTER_PASSWORD environment variable.")
		return fmt.Errorf("Missing ONEPASSWORD_MASTER_PASSWORD environment variable.")
	}

	err := b.OnePasswordClient.SignIn(domain, email, secretKey, masterPassword)

	os.Unsetenv("ONEPASSWORD_DOMAIN")
	os.Unsetenv("ONEPASSWORD_EMAIL")
	os.Unsetenv("ONEPASSWORD_SECRET_KEY")
	os.Unsetenv("ONEPASSWORD_MASTER_PASSWORD")

	if err != nil {
		return err
	}

	fmt.Println("Signed into 1password successfully.")

	return nil
}

// Retrieve the 1password item whose name matches the key and return the value of the 'password' field.
func (b *OnePasswordBackend) Get(key string) (string, error) {
	fmt.Println("Retrieving 1password item '" + key + "'.")

	item := b.OnePasswordClient.Get(key)
	if item == "" {
		return "", fmt.Errorf("Could not retrieve 1password item '" + key + "'.")
	}

	value := gjson.Get(item, "details.fields.#[name==\"password\"].value")
	if !value.Exists() {
		return "", fmt.Errorf("1password item '" + key + "' does not have a 'password' field.")
	}

	fmt.Println("1password item '" + key + "' value of 'password' field retrieved successfully.")

	return value.String(), nil
}

type OnePasswordClient interface {
	Get(key string) string
	SignIn(domain string, email string, secretKey string, masterPassword string) error
}

type OnePasswordCliClient struct {
}

func (c OnePasswordCliClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	fmt.Println("Signing into 1password.")

	cmd := exec.Command("/usr/local/bin/op", "signin", domain, email)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr

	b, err := pty.Start(cmd)
	if err != nil {
		fmt.Println(err, "/usr/local/bin/op signin failed with %s")
		return err
	}

	go func() {
		b.Write([]byte(secretKey + "\n"))
		b.Write([]byte{4})
		b.Write([]byte{4})
		b.Write([]byte(masterPassword + "\n"))
		b.Write([]byte{4})
		b.Write([]byte{4})
	}()

	fmt.Println("Started '/usr/local/bin/op'.")

	cmd.Wait()

	r, _ := regexp.Compile("export OP_SESSION_externalsecretoperator=\"(.+)\"")
	matches := r.FindAllStringSubmatch(outb.String(), -1)

	if len(matches) == 0 {
		fmt.Println("Could not retrieve token from 1password.")
		return nil
	}

	token := matches[0][1]
	fmt.Println("\nUpdated 'OP_SESSION_externalsecretoperator' environment variable.")
	os.Setenv("OP_SESSION_externalsecretoperator", token)

	return nil
}

// Invoke $ op get item 'key'
func (c OnePasswordCliClient) Get(key string) string {
	cmd := exec.Command("/usr/local/bin/op", "get", "item", key)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stderr.Bytes()))
		fmt.Println(string(stdout.Bytes()))
		fmt.Println(err, "/usr/local/bin/op get item '%s' failed: (%v)", key, err)
		return ""
	}
	return string(stdout.Bytes())
}
