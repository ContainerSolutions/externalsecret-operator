// Package akv implements backend for Azure Key Vault secrets
package akv

import (
	"context"
	"encoding/json"
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
	jsonErr := json.Unmarshal(credentials, &akvCred)
	if jsonErr != nil {
		log.Error(jsonErr, "")
		return jsonErr
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

	//Remove this after/if this issue gets fixed: https://github.com/Azure/azure-sdk-for-go/issues/13641
	if err := os.Setenv("AZURE_AUTH_LOCATION", file.Name()); err != nil {
		log.Error(err, "")
		return err
	}
	//end Remove

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

	secretResp, err := a.client.GetSecret(context.Background(), "https://"+a.keyvault+".vault.azure.net", key, version)
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	log.Info("Get secret succeeded")

	return *secretResp.Value, nil
}

//AzureCredentials needed to access the Key Vault service
type AzureCredentials struct {
	TennantID    string `json:"tennant_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Keyvault     string `json:"akvName"`
}
