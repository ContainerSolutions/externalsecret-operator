// Package credstash implements backend for Credstash (that uses AWS KMS and DynamoDB)
// Heavily inspired in github.com/ouzi-dev/credstash-operator using https://github.com/versent/unicreds
package credstash

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	"github.com/containersolutions/externalsecret-operator/pkg/utils"
	unicreds "github.com/versent/unicreds"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	defaultRegion          = "eu-west-2"
	credstashVersionLength = 19
)

var (
	log                     = ctrl.Log.WithName("credstash")
	table                   = ""
	configEncryptionContext = make(map[string]string)
)

// SecretManagerClient will be our unicreds client
type SecretManagerClient interface {
	SetKMSConfig(config *aws.Config)
	SetDynamoDBConfig(config *aws.Config)
	GetHighestVersionSecret(tableName *string, name string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error)
	GetSecret(tableName *string, name string, version string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error)
}

// SecretManagerClientStruct defining this struct to write methods for it
type SecretManagerClientStruct struct {
}

// SetKMSConfig sets configuration for KMS access
func (s SecretManagerClientStruct) SetKMSConfig(config *aws.Config) {
	unicreds.SetKMSConfig(config)
}

// SetDynamoDBConfig sets configuration for DynamoDB access
func (s SecretManagerClientStruct) SetDynamoDBConfig(config *aws.Config) {
	unicreds.SetDynamoDBConfig(config)
}

// GetHighestVersionSecret gets a secret with latest version from credstash
func (s SecretManagerClientStruct) GetHighestVersionSecret(tableName *string, name string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error) {
	return unicreds.GetHighestVersionSecret(tableName, name, encContext)
}

// GetSecret gets a secret with specific version from credstash
func (s SecretManagerClientStruct) GetSecret(tableName *string, name string, version string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error) {
	return unicreds.GetSecret(tableName, name, version, encContext)
}

// Backend represents a backend for Credstash
type Backend struct {
	SecretsManager SecretManagerClient
	session        *session.Session
}

func init() {
	backend.Register("credstash", NewBackend)
}

// NewBackend returns an uninitialized Backend for Credstash
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init initializes the Backend for Credstash
func (s *Backend) Init(parameters map[string]interface{}, credentials []byte) error {
	var err error

	s.session, err = utils.GetAWSSession(parameters, credentials, defaultRegion)
	if err != nil {
		return err
	}
	var ok bool
	table, ok = parameters["table"].(string)
	if !ok {
		log.Error(nil, "Credstash Dynamo DB table key missing")
		return nil
	}

	configEncryptionContext, ok = parameters["encryptionContext"].(map[string]string)
	if !ok {
		log.Info("Not using security encryption context. Consider using it")
	}

	s.SecretsManager = SecretManagerClientStruct{}
	s.SecretsManager.SetKMSConfig(s.session.Config)
	s.SecretsManager.SetDynamoDBConfig(s.session.Config)
	return nil
}

// Get retrieves the secret associated with key from Credstash
func (s *Backend) Get(key string, version string) (string, error) {
	if table == "" {
		table = "credential-store"
	}

	if s.SecretsManager == nil {
		log.Error(fmt.Errorf("error"), "backend not initialized")
		return "", fmt.Errorf("backend not initialized")
	}

	encryptionContext := unicreds.NewEncryptionContextValue()
	for k, v := range configEncryptionContext {
		if err := encryptionContext.Set(k + ":" + v); err != nil {
			return "", err
		}
	}
	if version == "" {
		creds, err := s.SecretsManager.GetHighestVersionSecret(aws.String(table), key, encryptionContext)
		if err != nil {
			log.Error(err, "Failed fetching secret from credstash",
				"Secret.Key", key, "Secret.Version", "latest", "Secret.Table", table, "Secret.Context", configEncryptionContext)

			return "", err
		}

		return creds.Secret, nil
	}
	formattedVersion, err := formatCredstashVersion(version)
	if err != nil {
		log.Error(err, "Failed formatting secret version",
			"Secret.Key", key, "Secret.Version", version, "Secret.Table", table, "Secret.Context", configEncryptionContext)
		return "", err
	}

	creds, err := s.SecretsManager.GetSecret(aws.String(table), key, formattedVersion, encryptionContext)
	if err != nil {
		log.Error(err, "Failed fetching secret from credstash",
			"Secret.Key", key, "Secret.Version", formattedVersion, "Secret.Table", table, "Secret.Context", configEncryptionContext)
		return "", err
	}

	return creds.Secret, nil
}

func formatCredstashVersion(inputVersion string) (string, error) {
	_, err := strconv.Atoi(inputVersion)
	if err != nil {
		log.Error(err, "Could not parse credstash version into number",
			"Secret.Version", inputVersion)
		return "", err
	}

	// we already have a padded version so nothing to do
	if len(inputVersion) == credstashVersionLength {
		return inputVersion, nil
	}

	// version is too longßß
	if len(inputVersion) > credstashVersionLength {
		return "", fmt.Errorf("version string is longer than supported. Maximum length is %d characters",
			credstashVersionLength)
	}

	// pad version with leading zeros until we reach credstashVersionLength
	// format becomes something like %019s which means pad the string until there's 19 0s
	format := fmt.Sprintf("%s%ds", "%0", credstashVersionLength)
	newVersion := fmt.Sprintf(format, inputVersion)

	return newVersion, nil
}
