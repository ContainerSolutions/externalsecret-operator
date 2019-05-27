package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
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

	fmt.Println("Signing in to 1password with email " + email + " and domain " + domain)

	err := b.OnePasswordClient.SignIn(domain, email, secretKey, masterPassword)

	os.Unsetenv("ONEPASSWORD_DOMAIN")
	os.Unsetenv("ONEPASSWORD_EMAIL")
	os.Unsetenv("ONEPASSWORD_SECRET_KEY")
	os.Unsetenv("ONEPASSWORD_MASTER_PASSWORD")

	if err != nil {
		return err
	}

	fmt.Println("Signed in to 1password successfully.")

	return nil
}

// Call the 1password client and parse the 'fields' array in the output. Return the 'v' property of the field object of which the 'n' property matches parameter key.
func (b *OnePasswordBackend) Get(key string) (string, error) {
	fmt.Println("Retrieving key " + key + " from 1password")

	opItemString := b.OnePasswordClient.Get(key)

	var opItem OpItem

	json.Unmarshal([]byte(opItemString), &opItem)

	var value = opItem.Details.Sections[0].Fields[0].V

	fmt.Println("Retrieved value from 1password")

	return value, nil
}

type OnePasswordClient interface {
	Get(key string) string
	SignIn(domain string, email string, secretKey string, masterPassword string) error
}

type OnePasswordCliClient struct {
}

func (c OnePasswordCliClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	fmt.Println("Signing into 1password via '/usr/local/bin/op'.")

	cmd := exec.Command("/usr/local/bin/op", "signin", domain, email)

	cmd.Stdout = os.Stdout
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
	io.Copy(os.Stdout, b)

	fmt.Println("Started '/usr/local/bin/op'.")

	cmd.Wait()

	fmt.Println("Signed into 1password successfully.")

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
	}
	return string(stdout.Bytes())
}

type OpItem struct {
	Details Details
}

type Details struct {
	Sections []Section
}

type Section struct {
	Fields []Field
}
type Field struct {
	N string
	V string
}
