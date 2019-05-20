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

func init() {
	BackendRegister("asm", NewAWSSecretsManagerBackend)
}

func NewAWSSecretsManagerBackend() BackendIface {
	backend := &AWSSecretsManagerBackend{}
	return backend
}

func (s *AWSSecretsManagerBackend) Init(params ...interface{}) error {
	var err error

	s.config, err = awsConfigFromParams(params...)
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

func awsConfigFromParams(params ...interface{}) (*aws.Config, error) {

	paramMap, err := paramsToMap(params...)
	if err != nil {
		return nil, err
	}

	accessKeyID := paramMap["accessKeyID"]
	secretAccessKey := paramMap["secretAccessKey"]
	region := paramMap["region"]

	return &aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			region),
	}, nil
}

func paramsToMap(params ...interface{}) (map[string]string, error) {

	paramKeys := []string{"accessKeyID", "secretAccessKey", "region"}

	if len(params) < 1 {
		return nil, fmt.Errorf("Invalid init parameters: not found %v", paramKeys)
	}

	paramType := reflect.TypeOf(params[0].(map[string]string))
	if paramType != reflect.TypeOf(map[string]string{}) {
		return nil, fmt.Errorf("Invalid init parameters: expected `map[string]string` found `%v", paramType)
	}

	paramMap := params[0].(map[string]string)

	for _, key := range paramKeys {
		paramValue, found := paramMap[key]
		if !found {
			return nil, fmt.Errorf("Invalid init paramters: expected `%v` not found", key)
		}

		paramType := reflect.TypeOf(paramValue)
		if paramType.Kind() != reflect.String {
			return nil, fmt.Errorf("Invalid init paramters: expected `%v` of type `string` got `%v`", key, paramType)
		}
	}

	return paramMap, nil
}
