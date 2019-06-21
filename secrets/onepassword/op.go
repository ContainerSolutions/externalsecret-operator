package onepassword

import (
	"fmt"
	"regexp"
)

type OP interface {
	SignIn(domain string, email string, secretKey string, masterPassword string) (Session, error)
}

type OPProcess struct {
	Command OPCommand
}

func (op *OPProcess) SignIn(domain string, email string, secretKey string, masterPassword string) (Session, error) {
	stdout, _ := op.Command.EnterCredentials(secretKey, masterPassword)

	fmt.Println(stdout)

	r, _ := regexp.Compile("export OP_SESSION_(.+)=\"(.+)\"")
	matches := r.FindAllStringSubmatch(stdout, -1)

	if len(matches) == 0 {
		return Session{"", ""}, fmt.Errorf("could not retrieve session from 1password")
	}

	fmt.Println(matches)

	key := matches[0][1]
	value := matches[0][2]

	return Session{key, value}, nil
}

type OPCommand interface {
	EnterCredentials(secretKey string, masterPassword string) (string, error)
}

type OPBinary struct {
}

func (oc *OPBinary) EnterCredentials(secretKey string, masterPassword string) (string, error) {
	return "", nil
}
