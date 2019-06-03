package secrets

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type MockOnePasswordClient struct {
	mock.Mock
}

func (m MockOnePasswordClient) SignIn(domain string, email string, secretKey string, masterPassword string) error {
	args := m.Called(domain, email, secretKey, masterPassword)
	return args.Error(0)
}

// Return a static JSON output for $ op get item 'testkey'
func (m MockOnePasswordClient) Get(key string) string {
	return `{
		"uuid": "r4qk25ahjrurehsejazi3tz57e",
		"templateUuid": "001",
		"trashed": "N",
		"createdAt": "2019-05-29T09:13:12Z",
		"updatedAt": "2019-05-29T10:53:49Z",
		"changerUuid": "GI2HNNU3OBEBDJUO6IOBA4EEOY",
		"itemVersion": 2,
		"vaultUuid": "s63lunnfg3pgoiuvq7bcl6taju",
		"details": {
		  "fields": [
			{
			  "designation": "username",
			  "name": "username",
			  "type": "T",
			  "value": ""
			},
			{
			  "designation": "password",
			  "name": "password",
			  "type": "P",
			  "value": "testvalue"
			}
		  ],
		  "notesPlain": "",
		  "sections": []
		},
		"overview": {
		  "URLs": [],
		  "ainfo": "",
		  "pbe": 0,
		  "pgrng": false,
		  "ps": 40,
		  "tags": [],
		  "title": "TestItem",
		  "url": ""
		}
	  }
	`
}

func TestGetOnePassword(t *testing.T) {
	secretKey := "testkey"
	secretValue := "testvalue"
	expectedValue := secretValue

	Convey("Given an OPERATOR_CONFIG env var", t, func() {
		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = &MockOnePasswordClient{}

		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}

func TestInitOnePassword(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {

		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		client := &MockOnePasswordClient{}
		client.On("SignIn", domain, email, secretKey, masterPassword).Return(nil)

		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"domain":         domain,
				"email":          email,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
			}

			backend.Init(params)

			Convey("Backend signs in via 1password client", func() {
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

		client := &MockOnePasswordClient{}
		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"domain":         domain,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "Error reading 1password backend parameters: Invalid init parameters: expected `email` not found.")
		})
	})
}

func TestInitOnePassword_MissingDomain(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
		masterPassword := "MasterPassword12346!"

		client := &MockOnePasswordClient{}
		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"email":          email,
				"secretKey":      secretKey,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "Error reading 1password backend parameters: Invalid init parameters: expected `domain` not found.")
		})
	})
}

func TestInitOnePassword_MissingSecretKey(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		masterPassword := "MasterPassword12346!"

		client := &MockOnePasswordClient{}
		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"email":          email,
				"domain":         domain,
				"masterPassword": masterPassword,
			}

			So(backend.Init(params).Error(), ShouldEqual, "Error reading 1password backend parameters: Invalid init parameters: expected `secretKey` not found.")
		})
	})
}

func TestInitOnePassword_MissingMasterPassword(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		domain := "https://externalsecretoperator.1password.com"
		email := "externalsecretoperator@example.com"
		secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"

		client := &MockOnePasswordClient{}
		backend := NewOnePasswordBackend()
		(backend).(*OnePasswordBackend).Client = client

		Convey("When initializing", func() {
			params := map[string]string{
				"email":     email,
				"domain":    domain,
				"secretKey": secretKey,
			}

			So(backend.Init(params).Error(), ShouldEqual, "Error reading 1password backend parameters: Invalid init parameters: expected `masterPassword` not found.")
		})
	})
}
