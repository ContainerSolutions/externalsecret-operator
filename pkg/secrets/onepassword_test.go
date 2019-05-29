package secrets

import (
	"os"
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

	Convey("Given an initialized OnePasswordBackend", t, func() {
		backend := NewOnePasswordBackend("Personal", MockOnePasswordClient{})

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

		os.Setenv("ONEPASSWORD_DOMAIN", domain)
		os.Setenv("ONEPASSWORD_EMAIL", email)
		os.Setenv("ONEPASSWORD_SECRET_KEY", secretKey)
		os.Setenv("ONEPASSWORD_MASTER_PASSWORD", masterPassword)

		client := MockOnePasswordClient{}

		client.On("SignIn", domain, email, secretKey, masterPassword).Return(nil)

		backend := NewOnePasswordBackend("Personal", client)

		Convey("When initializing", func() {
			backend.Init()
			Convey("Backend signs in via 1password client", func() {
				client.AssertExpectations(t)
				So(os.Getenv("ONEPASSWORD_DOMAIN"), ShouldEqual, "")
				So(os.Getenv("ONEPASSWORD_EMAIL"), ShouldEqual, "")
				So(os.Getenv("ONEPASSWORD_SECRET_KEY"), ShouldEqual, "")
				So(os.Getenv("ONEPASSWORD_MASTER_PASSWORD"), ShouldEqual, "")
			})
		})
	})
}

func TestInitOnePassword_MissingDomain(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {
		os.Unsetenv("ONEPASSWORD_DOMAIN")

		client := MockOnePasswordClient{}

		backend := NewOnePasswordBackend("Personal", client)

		Convey("When initializing", func() {
			So(backend.Init().Error(), ShouldEqual, "Missing ONEPASSWORD_DOMAIN environment variable.")
		})
	})
}

func TestInitOnePassword_MissingEmail(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {

		os.Setenv("ONEPASSWORD_DOMAIN", "https://externalsecretoperator.1password.com")

		os.Unsetenv("ONEPASSWORD_EMAIL")

		client := MockOnePasswordClient{}

		backend := NewOnePasswordBackend("Personal", client)

		Convey("When initializing", func() {
			So(backend.Init().Error(), ShouldEqual, "Missing ONEPASSWORD_EMAIL environment variable.")
		})
	})
}

func TestInitOnePassword_MissingSecretKey(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {

		os.Setenv("ONEPASSWORD_DOMAIN", "https://externalsecretoperator.1password.com")
		os.Setenv("ONEPASSWORD_EMAIL", "externalsecretoperator@example.com")

		os.Unsetenv("ONEPASSWORD_SECRET_KEY")

		client := MockOnePasswordClient{}

		backend := NewOnePasswordBackend("Personal", client)

		Convey("When initializing", func() {
			So(backend.Init().Error(), ShouldEqual, "Missing ONEPASSWORD_SECRET_KEY environment variable.")
		})
	})
}

func TestInitOnePassword_MissingMasterPassword(t *testing.T) {
	Convey("Given a OnePasswordBackend", t, func() {

		os.Setenv("ONEPASSWORD_DOMAIN", "https://externalsecretoperator.1password.com")
		os.Setenv("ONEPASSWORD_EMAIL", "externalsecretoperator@example.com")
		os.Setenv("ONEPASSWORD_SECRET_KEY", "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ")

		os.Unsetenv("ONEPASSWORD_MASTER_PASSWORD")

		client := MockOnePasswordClient{}

		backend := NewOnePasswordBackend("Personal", client)

		Convey("When initializing", func() {
			So(backend.Init().Error(), ShouldEqual, "Missing ONEPASSWORD_MASTER_PASSWORD environment variable.")
		})
	})
}
