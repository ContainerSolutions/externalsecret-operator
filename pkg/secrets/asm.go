package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// AWSSecretsManagerBackend represents a backend for AWS Secrets Manager
type AWSSecretsManagerBackend struct {
	SecretsManager secretsmanageriface.SecretsManagerAPI
	config         *aws.Config
	session        *session.Session
}

func init() {
	BackendRegister("asm", NewAWSSecretsManagerBackend)
}

// NewAWSSecretsManagerBackend returns an uninitialized AWSSecretsManagerBackend
func NewAWSSecretsManagerBackend() Backend {
	return &AWSSecretsManagerBackend{}
}

// Init initializes the AWSSecretsManagerBackend
func (s *AWSSecretsManagerBackend) Init(parameters map[string]string) error {
	var err error

	s.config, err = awsConfigFromParams(parameters)
	if err != nil {
		return err
	}

	s.session, err = session.NewSession(s.config)
	if err != nil {
		return err
	}

	s.SecretsManager = secretsmanager.New(s.session)
	return nil
}

// Get retrieves the secret associated with key from AWSSecretsManagerBackend
func (s *AWSSecretsManagerBackend) Get(key string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	err := input.Validate()
	if err != nil {
		return "", err
	}

	if s.SecretsManager == nil {
		return "", fmt.Errorf("backend not initialized")
	}

	output, err := s.SecretsManager.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *output.SecretString, nil
}

// awsConfigFromParams returns an aws.Config based on the parameters
func awsConfigFromParams(parameters map[string]string) (*aws.Config, error) {

	keys := []string{"accessKeyID", "secretAccessKey", "region"}

	for _, key := range keys {
		_, found := parameters[key]
		if !found {
			return nil, fmt.Errorf("Invalid init paramters: expected `%v` not found", key)
		}
	}

	return &aws.Config{
		Region: aws.String(parameters["region"]),
		Credentials: credentials.NewStaticCredentials(
			parameters["accessKeyID"],
			parameters["secretAccessKey"],
			""),
	}, nil
}
