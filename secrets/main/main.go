package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

func main() {
	opEmail := os.Getenv("OP_EMAIL")
	opDomain := os.Getenv("OP_DOMAIN")
	opSecretKey := os.Getenv("OP_SECRET_KEY")
	opMasterPassword := os.Getenv("OP_MASTER_PASSWORD")

	cmd, terminal := startOnePassword(opDomain, opEmail)
	defer func() { _ = terminal.Close() }()

	go func() {
		terminal.Write([]byte(opSecretKey + "\n"))
		terminal.Write([]byte{4})

		terminal.Write([]byte(opMasterPassword + "\n"))
		terminal.Write([]byte{4})
	}()

	err := cmd.Wait()
	if err != nil {
		fmt.Println(err, "/usr/local/bin/op signin failed with %s")
	}
}

func startOnePassword(domain string, email string) (cmd *exec.Cmd, terminal *os.File) {
	cmd = exec.Command("/usr/local/bin/op", "signin", domain, email)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	terminal, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Failed to start command", err)
	}

	return cmd, terminal
}
