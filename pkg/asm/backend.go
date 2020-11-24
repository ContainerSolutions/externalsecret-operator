// Package asm implements an external secret backend for AWS Secrets Manager.
package asm

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	"github.com/containersolutions/externalsecret-operator/pkg/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	defaultRegion = "eu-west-2"
)

var (
	log = ctrl.Log.WithName("asm")
)

// Backend represents a backend for AWS Secrets Manager
type Backend struct {
	SecretsManager secretsmanageriface.SecretsManagerAPI
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
func (s *Backend) Init(parameters map[string]interface{}, credentials []byte) error {
	var err error

	s.session, err = utils.GetAWSSession(parameters, credentials, defaultRegion)
	if err != nil {
		return err
	}

	s.SecretsManager = secretsmanager.New(s.session)
	return nil
}

// Get retrieves the secret associated with key from AWS Secrets Manager
func (s *Backend) Get(key string, version string) (string, error) {
	_ = version

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	err := input.Validate()
	if err != nil {
		return "", err
	}

	if s.SecretsManager == nil {
		log.Error(fmt.Errorf("error"), "backend not initialized")
		return "", fmt.Errorf("backend not initialized")
	}

	result, err := s.SecretsManager.GetSecretValue(input)
	if err != nil {
		log.Error(err, "Error getting secret value")
		return "", err
	}

	// https: //docs.aws.amazon.com/secretsmanager/latest/apireference/API_CreateSecret.html
	// TLDR: Either SecretString or SecretBinary must have a value, but not both. They cannot both be empty.
	var secretValue string
	if result.SecretString != nil {
		secretValue = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Error(err, "Base64 Decode Error:")
			return "", err
		}
		secretValue = string(decodedBinarySecretBytes[:len])
	}
	return secretValue, nil
}
