package onepassword

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type OnePasswordClient interface {
	Get(key string) string
	SignIn(domain string, email string, secretKey string, masterPassword string) error
}

type OnePasswordCliClient struct {
	Op OP
}

func (c OnePasswordCliClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	session, _ := c.Op.SignIn(domain, email, secretKey, masterPassword)

	os.Setenv(session.Key, session.Value)

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

// cmd := exec.Command("/usr/local/bin/op", "signin", domain, email)

// terminal, err := pty.Start(cmd)
// if err != nil {
// 	fmt.Println(err, "/usr/local/bin/op signin failed with %s")
// 	return err
// }
// defer func() { _ = terminal.Close() }()

// go func() {
// 	terminal.Write([]byte(secretKey + "\n"))
// 	terminal.Write([]byte{4})

// 	terminal.Write([]byte(masterPassword + "\n"))
// 	terminal.Write([]byte{4})
// }()

// fmt.Println("Started '/usr/local/bin/op'.")

// err = cmd.Wait()
// if err != nil {
// 	fmt.Println(err, "/usr/local/bin/op signin failed with %s")
// 	return err
// }

// r, _ := regexp.Compile("export OP_SESSION_(.+)=\"(.+)\"")
// matches := r.FindAllStringSubmatch(stdout.String(), -1)

// if len(matches) == 0 {
// 	fmt.Println("Could not retrieve token from 1password.")
// 	return nil
// }

// session := matches[0][1]
// token := matches[0][2]
// fmt.Println("\nUpdated 'OP_SESSION_externalsecretoperator' environment variable.")
// os.Setenv("OP_SESSION_"+session, token)
