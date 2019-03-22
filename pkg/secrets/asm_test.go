package secrets

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	. "github.com/smartystreets/goconvey/convey"
)

type mockedSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
}

func (m *mockedSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	mockSecretString := *input.SecretId + "Value"
	output := &secretsmanager.GetSecretValueOutput{
		Name:         input.SecretId,
		SecretString: &mockSecretString,
	}
	return output, nil
}

func TestGet(t *testing.T) {
	secretKey := "secret"
	secretValue := "secretValue"
	expectedValue := secretValue

	Convey("Given an initialized AWSSecretsManagerBackend", t, func() {
		backend := AWSSecretsManagerBackend{&mockedSecretsManager{}}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}
