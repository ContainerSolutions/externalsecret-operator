package onepassword

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

type Client interface {
	Get(key string) string
	SignIn(domain string, email string, secretKey string, masterPassword string) (Session, error)
}

type CliClient struct {
	Executable Executable
}

func (c CliClient) SignIn(domain string, email string, secretKey string, masterPassword string) (Session, error) {
	stdout, _ := c.Executable.SignIn(domain, email, secretKey, masterPassword)

	results, _ := regexp.Compile("export (OP_SESSION_(.+))=\"(.+)\"")
	matches := results.FindAllStringSubmatch(stdout, -1)

	if len(matches) == 0 {
		return Session{"", ""}, fmt.Errorf("could not retrieve session from 1password")
	}

	key := matches[0][1]
	value := matches[0][3]

	return Session{key, value}, nil
}

// Invoke $ op get item 'key'
func (c CliClient) Get(key string) string {
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
