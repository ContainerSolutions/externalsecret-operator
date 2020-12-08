// Package akv implements backend for Azure Key Vault secrets
package akv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("akv")

type ClientInterface interface {
	GetSecret(context context.Context, url string, key string, version string) (keyvault.SecretBundle, error)
}

// Backend represents a backend for Azure Key Vault
type Backend struct {
	Client   ClientInterface
	keyvault string
}

// NewBackend returns an uninitialized Backend for AWS Secret Manager
func NewBackend() backend.Backend {
	return &Backend{}
}

func init() {
	backend.Register("akv", NewBackend)
}

// Init initializes the Backend for Azure Key Vault
func (a *Backend) Init(parameters map[string]interface{}, credentials []byte) error {

	akvCred := AzureCredentials{}
	err := json.Unmarshal(credentials, &akvCred)
	if err != nil {
		log.Error(err, "")
		return err
	}

	file, err := ioutil.TempFile("/tmp", "akv")
	if err != nil {
		log.Error(err, "")
		return err
	}

	_, err = file.Write(credentials)
	if err != nil {
		log.Error(err, "error writing the temp file")
		return err
	}

	if err := os.Setenv("AZURE_AUTH_LOCATION", file.Name()); err != nil {
		log.Error(err, "error setting AZURE_AUTH_LOCATION environment variable")
		return err
	}

	authorizer, err := kvauth.NewAuthorizerFromFile(file.Name())
	if err != nil {
		log.Error(err, "error creating authorizer")
		return err
	}

	client := keyvault.New()
	client.Authorizer = authorizer
	a.Client = client
	a.keyvault = akvCred.Keyvault

	return nil
}

// Get retrieves the secret associated with key from Azure Key Vault
func (a *Backend) Get(key string, version string) (string, error) {

	if a.Client == nil {
		return "", errors.New("Azure Key Vault backend not initialized")
	}

	secretResp, err := a.Client.GetSecret(context.Background(), fmt.Sprintf("https://%s.vault.azure.net", a.keyvault), key, version)
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	log.Info("Get secret succeeded")

	return *secretResp.Value, nil
}

// AzureCredentials represents expected credentials
type AzureCredentials struct {
	TenantID     string `json:"tenantId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Keyvault     string `json:"keyvault"`
}
