package onepassword

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockCommand struct {
	mock.Mock
}

func (m *MockCommand) Start() error {
	m.Called()
	return nil
}

func (m *MockCommand) EnterCredentials(secretKey string, masterPassword string) (string, error) {
	return `export OP_SESSION_externalsecretoperator="gnKhaUdTwsmM1ESz4Q6cakpvtXMAEom7AAw04_xBi9s"
	# This command is meant to be used with your shell's eval function.
	# Run 'eval $(op signin cs)' to sign into your 1Password account.
	# If you wish to use the session token itself, pass the --output=raw flag value.`, nil
}

func TestOpSignIn(t *testing.T) {

	Convey("The session token is parsed after entering credentials", t, func() {

		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		command := &MockCommand{}

		op := OPProcess{}
		op.Command = command

		session, _ := op.SignIn(domain, email, secretKey, masterPassword)

		command.AssertExpectations(t)

		So(session, ShouldEqual, Session{"OP_SESSION_externalsecretoperator", "gXKhaUdTwsmM1ESz4Q6cakpvtXMAEom7AAw04_xB39s"})
	})
}
