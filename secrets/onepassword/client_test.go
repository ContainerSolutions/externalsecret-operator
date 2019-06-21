package onepassword

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockOp struct {
	mock.Mock
}

func (m *MockOp) SignIn(domain string, email string, secretKey string, masterPassword string) (Session, error) {
	return Session{"OP_SESSION_externalsecretoperator", "123456"}, nil
}

func TestSignIn(t *testing.T) {

	Convey("Given a OnePasswordCliClient", t, func() {
		op := &MockOp{}

		client := OnePasswordCliClient{}
		client.Op = op

		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		Convey("Session token should be set as an environment var after signing in", func() {
			client.SignIn(domain, email, secretKey, masterPassword)

			So(os.Getenv("OP_SESSION_externalsecretoperator"), ShouldEqual, "123456")

			op.AssertExpectations(t)
		})
	})
}
