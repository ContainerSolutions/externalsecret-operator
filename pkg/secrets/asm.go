package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type AWSSecretsManagerBackend struct {
	SecretsManager secretsmanageriface.SecretsManagerAPI
}

func NewAWSSecretsManagerBackend() *AWSSecretsManagerBackend {
	backend := &AWSSecretsManagerBackend{}
	backend.Init()
	return backend
}

func (s *AWSSecretsManagerBackend) Init(params ...interface{}) error {
	session, err := session.NewSession()
	if err != nil {
		return err
	}
	_, err = session.Config.Credentials.Get()
	if err != nil {
		return err
	}
	s.SecretsManager = secretsmanager.New(session)
	return nil
}

func (s *AWSSecretsManagerBackend) Get(key string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}
	err := input.Validate()
	if err != nil {
		return "", err
	}

	output, err := s.SecretsManager.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *output.SecretString, nil
}
