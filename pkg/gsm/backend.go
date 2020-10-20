// Package gsm implements backend for Google Secrets Manager
package gsm

import (
	"context"
	"fmt"

	"cloud.google.com/go/iam"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	"github.com/googleapis/gax-go"
	"golang.org/x/oauth2/google"
	option "google.golang.org/api/option"
	"google.golang.org/grpc"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	// iampb "google.golang.org/genproto/googleapis/iam/v1"
)

const (
	cloudPlatformRole = "https://www.googleapis.com/auth/cloud-platform"
)

var log = ctrl.Log.WithName("gsm")

type GoogleSecretManagerClient interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
	AddSecretVersion(ctx context.Context, req *secretmanagerpb.AddSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	Connection() *grpc.ClientConn
	CreateSecret(ctx context.Context, req *secretmanagerpb.CreateSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error)
	DeleteSecret(ctx context.Context, req *secretmanagerpb.DeleteSecretRequest, opts ...gax.CallOption) error
	DestroySecretVersion(ctx context.Context, req *secretmanagerpb.DestroySecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	DisableSecretVersion(ctx context.Context, req *secretmanagerpb.DisableSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	EnableSecretVersion(ctx context.Context, req *secretmanagerpb.EnableSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error)
	GetSecret(ctx context.Context, req *secretmanagerpb.GetSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error)
	GetSecretVersion(ctx context.Context, req *secretmanagerpb.GetSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.SecretVersion, error)
	IAM(name string) *iam.Handle
	ListSecretVersions(ctx context.Context, req *secretmanagerpb.ListSecretVersionsRequest, opts ...gax.CallOption) *secretmanager.SecretVersionIterator
	ListSecrets(ctx context.Context, req *secretmanagerpb.ListSecretsRequest, opts ...gax.CallOption) *secretmanager.SecretIterator
	SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error)
	TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error)
	UpdateSecret(ctx context.Context, req *secretmanagerpb.UpdateSecretRequest, opts ...gax.CallOption) (*secretmanagerpb.Secret, error)
	Close() error
}

// Backend for Google Secrets Manager
type Backend struct {
	projectID           string
	SecretManagerClient GoogleSecretManagerClient
}

func init() {
	backend.Register("gsm", NewBackend)
}

// NewBackend gives you an empty Google Secrets Manager Backend
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init initializes Google secretsmanager backend
func (g *Backend) Init(parameters map[string]interface{}, credentials []byte) error {
	ctx := context.Background()

	if len(parameters) == 0 || len(credentials) == 0 {
		return fmt.Errorf("credentials or parameters invalid")
	}

	projectID, ok := parameters["projectID"].(string)
	if !ok {
		return fmt.Errorf("parameters invalid")
	}

	g.projectID = projectID

	config, err := google.JWTConfigFromJSON(credentials, cloudPlatformRole)
	if err != nil {
		return err
	}

	ts := config.TokenSource(ctx)

	client, err := secretmanager.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	g.SecretManagerClient = client

	return nil
}

// Get a key and returns a value
func (g *Backend) Get(key string, version string) (string, error) {
	ctx := context.Background()

	if g.SecretManagerClient == nil {
		return "", fmt.Errorf("backend is not initialized")
	}

	validVersion := version
	if validVersion == "" {
		validVersion = "latest"
	}

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", g.projectID, key, validVersion)

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := g.SecretManagerClient.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}
