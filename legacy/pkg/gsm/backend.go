// Package gsm implements backend for Google Secrets Manager
package gsm

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	"golang.org/x/oauth2/google"
	option "google.golang.org/api/option"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	cloudPlatformRole = "https://www.googleapis.com/auth/cloud-platform"
)

var log = logf.Log.WithName("gsm")

// Backend for Google Secrets Manager
type Backend struct {
	projectID           string
	SecretManagerClient *secretmanager.Client
}

func init() {
	backend.Register("gsm", NewBackend)
}

// NewBackend gives you an empty Google Secrets Manager Backend
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init initializes Google secretsmanager backend
func (g *Backend) Init(parameters map[string]string) error {
	ctx := context.Background()
	g.projectID = parameters["projectID"]

	sAccount := serviceAccount{}
	jsonCredentials, err := sAccount.Marshal(parameters)
	println(string(jsonCredentials))
	if err != nil {
		log.Info("here")
		return err
	}

	config, err := google.JWTConfigFromJSON(jsonCredentials, cloudPlatformRole)
	if err != nil {
		log.Info("here")
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
func (g *Backend) Get(key string) (string, error) {
	ctx := context.Background()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", g.projectID, key, "latest")

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := g.SecretManagerClient.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %v", err)
	}

	// fmt.Fprintf(os.Stdout, "Plaintext: %s\n", string(result.Payload.Data))

	return string(result.Payload.Data), nil
}

type serviceAccount struct {
	AuthType                string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func (s *serviceAccount) Marshal(param map[string]string) ([]byte, error) {
	s.AuthType = param["type"]
	s.ProjectID = param["projectID"]
	s.PrivateKeyID = param["privateKeyID"]
	s.PrivateKey = param["privateKey"]
	s.ClientEmail = param["clientEmail"]
	s.ClientID = param["clientID"]
	s.AuthURI = param["authURI"]
	s.TokenURI = param["tokenURI"]
	s.AuthProviderX509CertURL = param["authProviderX509CertURL"]
	s.ClientX509CertURL = param["clientX509CertURL"]

	return json.Marshal(s)
}
