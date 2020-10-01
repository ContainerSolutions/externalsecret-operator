package gsm

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// type mockGoogleClient struct {
// 	// connPool         gtransport.ConnPool
// 	// disableDeadlines bool
// 	// client           secretmanagerpb.SecretManagerServiceClient
// 	// CallOptions      *secretmanager.CallOptions
// 	// xGoogMetadata    metadata.MD
// 	client *secretmanager.Client
// }

// type mockGoogleClientInterface interface {
// 	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error)
// }

// func (g *mockGoogleClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
// 	return &secretmanagerpb.AccessSecretVersionResponse{
// 		Name: "test",
// 		Payload: &secretmanagerpb.SecretPayload{
// 			Data: []byte("Testing"),
// 		},
// 	}, nil
// }

func TestNewBackend(t *testing.T) {
	Convey("When creating a new GSM backend", t, func() {
		backend := NewBackend()
		So(backend, ShouldNotBeNil)
		So(backend, ShouldHaveSameTypeAs, &Backend{})
	})
}

func TestGet(t *testing.T) {
	secretKey := "secret"
	keyVersion := ""

	Convey("Given an uninitialized GoogleSecretsManager", t, func() {
		backend := Backend{}
		Convey("When retrieving a secret", func() {
			_, err := backend.Get(secretKey, keyVersion)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "backend is not initialized")
			})
		})
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
