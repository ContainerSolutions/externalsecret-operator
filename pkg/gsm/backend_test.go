package gsm

import (
	"context"
	"fmt"
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
	secretName := req.Name
	if secretName == "projects/test-project-gsm/secrets/SecretKeyError/versions/latest" {
		return nil, fmt.Errorf("Mocked errror")
	}

	return &secretmanagerpb.AccessSecretVersionResponse{
		Name: secretName,
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(secretName),
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
	var (
		secretKey      = "SecretKey"
		secretKeyError = "SecretKeyError"
		keyVersion     = "latest"
		testProject    = "test-project-gsm"
	)

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
		backend.projectID = testProject
		backend.SecretManagerClient = &mockGoogleSecretManagerClient{}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, "projects/test-project-gsm/secrets/SecretKey/versions/latest")
			})
		})
	})

	Convey("Given an initialized GoogleSecretManger Client", t, func() {
		backend := Backend{}
		backend.projectID = testProject
		backend.SecretManagerClient = &mockGoogleSecretManagerClient{}
		Convey("When retrieving a secret with a empty version", func() {
			actualValue, err := backend.Get(secretKey, "")
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, "projects/test-project-gsm/secrets/SecretKey/versions/latest")
			})
		})
	})

	Convey("Given an initialized GoogleSecretManger Client", t, func() {
		backend := Backend{}
		backend.projectID = testProject
		backend.SecretManagerClient = &mockGoogleSecretManagerClient{}
		Convey("When GetSecretValue() fails", func() {
			_, err := backend.Get(secretKeyError, "")
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "failed to access secret version: Mocked errror")
			})
		})
	})

}

func TestInit(t *testing.T) {

	Convey("During initilization", t, func() {
		var (
			backend     = Backend{}
			params      = make(map[string]interface{})
			credentials = make([]byte, 1, 1)
		)

		Convey("When parameters are blank", func() {
			err := backend.Init(params, credentials)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "credentials or parameters invalid")
			})
		})
	})

	Convey("During initilization", t, func() {
		var (
			backend     = Backend{}
			params      = make(map[string]interface{})
			credentials = make([]byte, 1, 1)
		)

		params["invalid"] = "invalid value"

		Convey("When parameters are invalid", func() {
			err := backend.Init(params, credentials)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "parameters invalid")
			})
		})
	})

	Convey("During initilization", t, func() {
		var (
			backend = Backend{}
			params  = make(map[string]interface{})
		)

		params["projectID"] = "test-project"

		Convey("When service account values are blank", func() {
			err := backend.Init(params, []byte{})
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "credentials or parameters invalid")
			})
		})
	})

	Convey("During initilization", t, func() {
		var (
			serviceAccount = `{
				"project_id": "external-secrets-operator",
				"private_key_id": ""
			}`
			backend = Backend{}
			params  = make(map[string]interface{})
		)

		params["projectID"] = "test-project"

		Convey("When a service account is invalid", func() {
			err := backend.Init(params, []byte(serviceAccount))
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "google: read JWT from JSON credentials")
			})
		})
	})

	Convey("During initilization", t, func() {
		var (
			serviceAccount = `{
				"type": "service_account",
				"project_id": "test-project",
				"private_key_id": "",
				"private_key": "-----BEGIN PRIVATE KEY-----\nA KEy\n-----END PRIVATE KEY-----\n",
				"client_email": "test-service-account@test-project.iam.gserviceaccount.com",
				"client_id": "",
				"auth_uri": "https://accounts.google.com/o/oauth2/auth",
				"token_uri": "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/operator-service-account%40external-secrets-operator.iam.gserviceaccount.com"
			}`
			backend    = Backend{}
			testClient = secretmanager.Client{}
			params     = make(map[string]interface{})
		)

		params["projectID"] = "test-project"

		Convey("When a service account is valid", func() {
			err := backend.Init(params, []byte(serviceAccount))
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(backend.SecretManagerClient, ShouldNotBeNil)
				So(backend.SecretManagerClient, ShouldHaveSameTypeAs, &testClient)
			})
		})
	})

}
