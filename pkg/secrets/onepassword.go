package secrets

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"

	"github.com/kr/pty"
	"github.com/tidwall/gjson"
)

type OnePasswordBackend struct {
	Client OnePasswordClient
	Vault  string
}

func NewOnePasswordBackend() Backend {
	backend := &OnePasswordBackend{}
	backend.Client = OnePasswordCliClient{}
	backend.Vault = "Personal"
	return backend
}

func init() {
	BackendRegister("onepassword", NewOnePasswordBackend)
}

// Read secrets from the parameters and sign in to 1password.
func (b *OnePasswordBackend) Init(parameters map[string]string) error {
	err := validateParameters(parameters)
	if err != nil {
		return fmt.Errorf("Error reading 1password backend parameters: %v", err)
	}

	b.Vault = parameters["vault"]

	err = b.Client.SignIn(parameters["domain"], parameters["email"], parameters["secretKey"], parameters["masterPassword"])
	if err != nil {
		return err
	}

	fmt.Println("Signed into 1password successfully.")

	return nil
}

// Retrieve the 1password item whose name matches the key and return the value of the 'password' field.
func (b *OnePasswordBackend) Get(key string) (string, error) {
	fmt.Println("Retrieving 1password item '" + key + "'.")

	item := b.Client.Get(key)
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

func validateParameters(parameters map[string]string) error {

	paramKeys := []string{"domain", "email", "secretKey", "masterPassword", "vault"}

	for _, key := range paramKeys {
		paramValue, found := parameters[key]
		if !found {
			return fmt.Errorf("Invalid init parameters: expected `%v` not found.", key)
		}

		paramType := reflect.TypeOf(paramValue)
		if paramType.Kind() != reflect.String {
			return fmt.Errorf("Invalid init parameters: expected `%v` of type `string` got `%v`", key, paramType)
		}
	}

	return nil
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
