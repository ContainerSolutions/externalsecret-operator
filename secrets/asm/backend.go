// Package asm implements an external secret backend for AWS Secrets Manager.
package asm

import (
	"fmt"

	"github.com/containersolutions/externalsecret-operator/secrets/backend"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// Backend represents a backend for AWS Secrets Manager
type Backend struct {
	SecretsManager secretsmanageriface.SecretsManagerAPI
	config         *aws.Config
	session        *session.Session
}

func init() {
	backend.Register("asm", NewBackend)
}

// NewBackend returns an uninitialized Backend for AWS Secret Manager
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init initializes the Backend for AWS Secret Manager
func (s *Backend) Init(parameters map[string]string) error {
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

// Get retrieves the secret associated with key from AWS Secrets Manager
func (s *Backend) Get(key string) (string, error) {
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
			return nil, fmt.Errorf("Invalid init parameters: expected `%v` not found", key)
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
