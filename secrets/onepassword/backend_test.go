package onepassword

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m MockClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	args := m.Called(domain, email, secretKey, masterPassword)
	return args.Error(0)
}

// Return a static JSON output for $ op get item 'testkey'
func (m MockClient) Get(value string, key string) (string, error) {
	return "testvalue", nil
}

func TestGetOnePassword(t *testing.T) {
	secretKey := "testkey"
	secretValue := "testvalue"
	expectedValue := secretValue

	Convey("Given an OPERATOR_CONFIG env var", t, func() {
		backend := NewBackend()
		(backend).(*Backend).Client = &MockClient{}

		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}

func TestOnePasswordBackend_DefaultVault(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		backend := NewBackend()

		Convey("The default vault should be 'Personal'", func() {
			So((backend).(*Backend).Vault, ShouldEqual, "Personal")
		})
	})
}

func TestInitOnePassword(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {

		vault := "production"

		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		client := &MockClient{}
		client.On("SignIn", domain, email, secretKey, masterPassword).Return(nil)

		backend := NewBackend()
		(backend).(*Backend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"domain":         domain,
				"email":          email,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
				"vault":          vault,
			}

			backend.Init(params)

			Convey("Client should have signed in", func() {
				client.AssertExpectations(t)
			})
		})
	})
}

func TestInitOnePassword_MissingEmail(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		domain := "https://externalsecretoperator.1password.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		backend := NewBackend()

		Convey("When initializing", func() {
			params := map[string]string{
				"domain":         domain,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "error reading 1password backend parameters: invalid init parameters: expected `email` not found")
		})
	})
}

func TestInitOnePassword_MissingDomain(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		backend := NewBackend()

		Convey("When initializing", func() {
			params := map[string]string{
				"email":          email,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "error reading 1password backend parameters: invalid init parameters: expected `domain` not found")
		})
	})
}

func TestInitOnePassword_MissingSecretKey(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		masterPassword := "MasterPassword12346!"

		backend := NewBackend()

		Convey("When initializing", func() {
			params := map[string]string{
				"email":          email,
				"domain":         domain,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "error reading 1password backend parameters: invalid init parameters: expected `secretKey` not found")
		})
	})
}

func TestInitOnePassword_MissingMasterPassword(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"

		backend := NewBackend()

		Convey("When initializing", func() {
			params := map[string]string{
				"email":     email,
				"domain":    domain,
				"secretKey": secretKey,
			}

			So(backend.Init(params).Error(), ShouldEqual, "error reading 1password backend parameters: invalid init parameters: expected `masterPassword` not found")
		})
	})
}
