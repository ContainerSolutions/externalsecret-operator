package secrets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	url := os.Getenv("ONEPASSWORD_DOMAIN")
	email := os.Getenv("ONEPASSWORD_EMAIL")
	secretKey := os.Getenv("ONEPASSWORD_SECRET_KEY")
	masterPassword := os.Getenv("ONEPASSWORD_MASTER_PASSWORD")

	if url == "" || email == "" || secretKey == "" || masterPassword == "" {
		return fmt.Errorf("Missing one or more ONEPASSWORD environment variables.")
	}

	err := b.OnePasswordClient.SignIn(url, email, secretKey, masterPassword)

	os.Unsetenv("ONEPASSWORD_DOMAIN")
	os.Unsetenv("ONEPASSWORD_EMAIL")
	os.Unsetenv("ONEPASSWORD_SECRET_KEY")
	os.Unsetenv("ONEPASSWORD_MASTER_PASSWORD")

	if err != nil {
		return err
	}
	return nil
}

// Call the 1password client and parse the 'fields' array in the output. Return the 'v' property of the field object of which the 'n' property matches parameter key.
func (b *OnePasswordBackend) Get(key string) (string, error) {
	opItemString := b.OnePasswordClient.Get(key)

	var opItem OpItem

	json.Unmarshal([]byte(opItemString), &opItem)

	var value = opItem.Details.Sections[0].Fields[0].V

	return value, nil
}

type OnePasswordClient interface {
	Get(key string) string
	SignIn(domain string, email string, secretKey string, masterPassword string) error
}

type OnePasswordCliClient struct {
}

func (c OnePasswordCliClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	log.Println("Signing into 1password via '/usr/local/bin/op'.")

	cmd := exec.Command("op", "signin", domain, email)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	b, err := pty.Start(cmd)
	if err != nil {
		log.Fatalf("/usr/local/bin/op signin failed with %s.\n", err)
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

	log.Println("Started '/usr/local/bin/op'.")

	cmd.Wait()

	log.Println("Signed into 1password successfully.")

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
		log.Fatalf("/usr/local/bin/op get item '"+key+"' failed with %s\n", err)
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
