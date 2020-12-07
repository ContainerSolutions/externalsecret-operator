// Package akv implements backend for Azure Key Vault secrets
package akv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("akv")

// Backend is the needed structure to access the Azure Key Vault service
type Backend struct {
	client   keyvault.BaseClient
	keyvault string
}

// NewBackend gives you a new akv.Backend
func NewBackend() backend.Backend {
	return &Backend{}
}

func init() {
	backend.Register("akv", NewBackend)
}

// Init will initialize and authorize a local client responsible for accessing the Azure Key Vault service
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

	a.client = keyvault.New()
	a.client.Authorizer = authorizer
	a.keyvault = akvCred.Keyvault

	return nil
}

// Get is responsible for getting the actual secret value from Azure Key Vault
func (a *Backend) Get(key string, version string) (string, error) {

	secretResp, err := a.client.GetSecret(context.Background(), fmt.Sprintf("https://%s.vault.azure.net", a.keyvault), key, version)
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	log.Info("Get secret succeeded")

	return *secretResp.Value, nil
}

//AzureCredentials needed to access the Key Vault service
type AzureCredentials struct {
	TenantID     string `json:"tenantId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Keyvault     string `json:"keyvault"`
}
