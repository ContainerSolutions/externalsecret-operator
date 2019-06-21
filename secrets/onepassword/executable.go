package onepassword

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/kr/pty"
)

type Executable interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) (string, error)
}

type OP struct {
}

func (o *OP) SignIn(domain string, email string, secretKey string, masterPassword string) (string, error) {
	cmd := exec.Command("/usr/local/bin/op", "signin", domain, email)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	terminal, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("failed to start 'op signin'", err)
		return "", err
	}

	defer func() { _ = terminal.Close() }()

	go func() {
		terminal.Write([]byte(secretKey + "\n"))
		terminal.Write([]byte{4})

		terminal.Write([]byte(masterPassword + "\n"))
		terminal.Write([]byte{4})
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err, "command 'op signin' failed with %s")
		return "", err
	}

	return stdout.String(), nil
}
