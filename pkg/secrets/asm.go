package secrets

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type AWSSecretsManagerBackend struct {
	Backend
	SecretsManager secretsmanageriface.SecretsManagerAPI
	config         *aws.Config
	session        *session.Session
}

func NewAWSSecretsManagerBackend() *AWSSecretsManagerBackend {
	backend := &AWSSecretsManagerBackend{}
	backend.Init()
	return backend
}

func (s *AWSSecretsManagerBackend) Init(params ...interface{}) error {
	var err error

	s.config, err = AWSConfigFromParams(params)
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

func AWSConfigFromParams(params ...interface{}) (*aws.Config, error) {
	nParams := 3

	if len(params) < nParams {
		return nil, fmt.Errorf("Invalid init paramters: aws_access_key_id, aws_secret_access_key, region")
	}

	for i := 0; i < nParams; i++ {
		t := reflect.TypeOf(params[i])
		if t.Kind() != reflect.String {
			return nil, fmt.Errorf("Invalid init paramters: expected `string` got `%v`", t)
		}
	}

	accessKeyID := params[0].(string)
	secretAccessKey := params[1].(string)
	region := params[2].(string)

	return &aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			region),
	}, nil
}
