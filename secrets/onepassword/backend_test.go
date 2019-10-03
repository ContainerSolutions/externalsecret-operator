package onepassword

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockCli struct {
	value string
}

func (m *MockCli) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	return nil
}

func (m *MockCli) GetItem(vault string, item string) (string, error) {
	if m.value != "" {
		return m.value, nil
	} else {
		return "", fmt.Errorf("invalid item")
	}
}

func TestGet(t *testing.T) {
	item := "item"
	value := "value"

	backend := &Backend{}
	backend.Vault = "Shared"
	backend.OnePassword = &MockCli{value: value}

	actual, err := backend.Get(item)
	if err != nil {
		t.Fail()
		fmt.Printf("expected nil but got '%s'", err)
	}
	if actual != value {
		t.Fail()
		fmt.Printf("expected '%s' got %s'", value, actual)
	}
}

func TestGet_ErrGetItem(t *testing.T) {
	backend := &Backend{}
	backend.OnePassword = &MockCli{}

	_, err := backend.Get("nonExistentItem")
	switch err.(type) {
	case *ErrGetItem:
	default:
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"
	vault := "production"

	backend := &Backend{
		OnePassword: &Cli{Op: &FakeOp{
			signInOk: true,
		}},
	}

	params := map[string]string{
		"domain":         domain,
		"email":          email,
		"secretKey":      secretKey,
		"masterPassword": masterPassword,
		"vault":          vault,
	}

	err := backend.Init(params)
	if err != nil {
		t.Fail()
		fmt.Println("expected signin to succeed")
	}
}

func TestInit_ErrSigninFailed(t *testing.T) {
	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"
	vault := "production"

	backend := &Backend{
		OnePassword: &Cli{Op: &FakeOp{
			signInOk: false,
		}},
	}

	params := map[string]string{
		"domain":         domain,
		"email":          email,
		"secretKey":      secretKey,
		"masterPassword": masterPassword,
		"vault":          vault,
	}

	err := backend.Init(params)
	switch err.(type) {
	case *ErrSigninFailed:
		actual := err.Error()
		expected := "could not sign in to 1password: fake op sign in programmed to fail"
		if actual != expected {
			t.Fail()
			fmt.Printf("expected '%s' got '%s'", expected, actual)
		}
	default:
		t.Fail()
		fmt.Println("expected signin failed error")
	}
}

func TestInit_ErrParameterMissing_Email(t *testing.T) {
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

func TestInit_ErrParameterMissing_Domain(t *testing.T) {
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

func TestInit_ErrParameterMissing_SecretKey(t *testing.T) {
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

func TestInit_ParameterMissing_MasterPassword(t *testing.T) {
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

func TestNewBackend(t *testing.T) {
	backend := NewBackend()

	if backend.(*Backend).OnePassword == nil {
		t.Fail()
		fmt.Println("expected backend to have a 1password cli")
	}

	expectedVault := "Personal"

	if backend.(*Backend).Vault != expectedVault {
		t.Fail()
		fmt.Printf("expected vault to be equal to '%s'", expectedVault)
	}
}
