package onepassword

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/kr/pty"
)

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
