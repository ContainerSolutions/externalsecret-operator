package gsm

import (
	"context"
	"testing"

	"cloud.google.com/go/iam"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/googleapis/gax-go"
	. "github.com/smartystreets/goconvey/convey"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc"
	// iampb "google.golang.org/genproto/googleapis/iam/v1"
)

type mockGoogleSecretManagerClient struct{}

func (g *mockGoogleSecretManagerClient) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return &secretmanagerpb.AccessSecretVersionResponse{
		Name: "test",
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte("Testing"),
		},
	}, nil
}

func (g *mockGoogleSecretManagerClient) Close() error {
	return nil
}

func (g *mockGoogleSecretManagerClient) AddSecretVersion(ctx context.Context, req *secretmanagerpb.AddSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) Connection() *grpc.ClientConn {
	return nil
}

func (g *mockGoogleSecretManagerClient) CreateSecret(ctx context.Context, req *secretmanagerpb.CreateSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) DeleteSecret(ctx context.Context, req *secretmanagerpb.DeleteSecretRequest, opts ...gax.CallOption) error {
	return nil
}

func (g *mockGoogleSecretManagerClient) DestroySecretVersion(ctx context.Context, req *secretmanagerpb.DestroySecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) DisableSecretVersion(ctx context.Context, req *secretmanagerpb.DisableSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) EnableSecretVersion(ctx context.Context, req *secretmanagerpb.EnableSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) GetSecret(ctx context.Context, req *secretmanagerpb.GetSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) GetSecretVersion(ctx context.Context, req *secretmanagerpb.GetSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error) {
	return nil, nil
}

func (g *mockGoogleSecretManagerClient) IAM(name string) *iam.Handle {
	return nil
}

func (g *mockGoogleSecretManagerClient) ListSecretVersions(ctx context.Context, req *secretmanagerpb.ListSecretVersionsRequest, opts ...gax.CallOption) *secretmanager.SecretVersionIterator {
	return nil
}

func (g *mockGoogleSecretManagerClient) ListSecrets(ctx context.Context, req *secretmanagerpb.ListSecretsRequest, opts ...gax.CallOption) *secretmanager.SecretIterator {
	return nil
}

func (g *mockGoogleSecretManagerClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return nil, nil
}
func (g *mockGoogleSecretManagerClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	return nil, nil
}
func (g *mockGoogleSecretManagerClient) UpdateSecret(ctx context.Context, req *secretmanagerpb.UpdateSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error) {
	return nil, nil
}

func TestNewBackend(t *testing.T) {
	Convey("When creating a new GSM backend", t, func() {
		backend := NewBackend()
		So(backend, ShouldNotBeNil)
		So(backend, ShouldHaveSameTypeAs, &Backend{})
	})
}

func TestGet(t *testing.T) {
	secretKey := "secret"
	keyVersion := "latest"

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

	Convey("Given an initialized GoogleSecretManger Client", t, func() {
		backend := Backend{}
		backend.SecretManagerClient = &mockGoogleSecretManagerClient{}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, "Testing")
			})
		})
	})

	Convey("Given an initialized GoogleSecretManger Client", t, func() {
		backend := Backend{}
		backend.SecretManagerClient = &mockGoogleSecretManagerClient{}
		Convey("When retrieving a secret with a empty version", func() {
			actualValue, err := backend.Get(secretKey, "")
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, "Testing")
			})
		})
	})
}

func TestInit(t *testing.T) {

	Convey("During initilization", t, func() {
		backend := Backend{}
		params := make(map[string]string)

		Convey("When parameters are blank", func() {
			err := backend.Init(params)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid or empty Config")
			})
		})
	})

	Convey("During initilization", t, func() {
		backend := Backend{}
		params := make(map[string]string)
		params["invalid"] = "invalid value"

		Convey("When parameters are invalid", func() {
			err := backend.Init(params)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid parameters")
			})
		})
	})

	Convey("During initilization", t, func() {
		backend := Backend{}
		params := make(map[string]string)
		params["projectID"] = "test-project"

		Convey("When service account values are blank or invalid", func() {
			err := backend.Init(params)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "google: read JWT from JSON credentials:")
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
