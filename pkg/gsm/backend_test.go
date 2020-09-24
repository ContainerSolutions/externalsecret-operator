package gsm

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBackend(t *testing.T) {
	Convey("When creating a new GSM backend", t, func() {
		backend := NewBackend()
		So(backend, ShouldNotBeNil)
		So(backend, ShouldHaveSameTypeAs, &Backend{})
	})
}

func TestSeviceAccountMarshal(t *testing.T) {
	Convey("When creating marshalling a service account", t, func() {
		params := map[string]string{
			"projectID":               "external-secrets-operator",
			"type":                    "service_account",
			"privateKeyID":            "pid",
			"privateKey":              "-----BEGIN PRIVATE KEY-----\nsome-key----END PRIVATE KEY-----\n",
			"clientEmail":             "operator-service-account@externalsecrets-operator.iam.gserviceaccount.com",
			"clientID":                "0000505056969",
			"authURI":                 "https://accounts.google.com/o/oauth2/auth",
			"tokenURI":                "https://oauth2.googleapis.com/token",
			"authProviderX509CertURL": "https://www.googleapis.com/oauth2/v1/certs",
			"clientX509CertURL":       "https://www.googleapis.com/robot/v1/metadata/x509/operator-service-account%40externalsecrets-operator.iam.gserviceaccount.com",
		}

		sAccount := serviceAccount{}
		jsonCredentials, _ := sAccount.Marshal(params)

		So(jsonCredentials, ShouldNotBeNil)

	})
}
