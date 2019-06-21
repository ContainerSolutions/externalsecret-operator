package onepassword

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockExecutable struct {
	mock.Mock
}

func (m *MockExecutable) SignIn(domain string, email string, secretKey string, masterPassword string) (string, error) {
	return `export OP_SESSION_externalsecretoperator="gXKhaUdTwsmM1ESz4Q6cakpvtXMAEom7AAw04_xB39s"
	# This command is meant to be used with your shell's eval function.
	# Run 'eval $(op signin cs)' to sign into your 1Password account.
	# If you wish to use the session token itself, pass the --output=raw flag value.`, nil
}

func TestSignIn(t *testing.T) {

	Convey("The session token is parsed after entering credentials", t, func() {

		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		executable := &MockExecutable{}

		client := CliClient{}
		client.Executable = executable

		session, _ := client.SignIn(domain, email, secretKey, masterPassword)

		executable.AssertExpectations(t)

		So(session.Key, ShouldEqual, "OP_SESSION_externalsecretoperator")
		So(session.Value, ShouldEqual, "gXKhaUdTwsmM1ESz4Q6cakpvtXMAEom7AAw04_xB39s")
	})
}
