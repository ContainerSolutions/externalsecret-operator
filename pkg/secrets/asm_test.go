package secrets

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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
		backend := AWSSecretsManagerBackend{}
		backend.SecretsManager = &mockedSecretsManager{}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}

func TestAWSConfigFromParams(t *testing.T) {

	Convey("Given a set of params", t, func() {
		expectedAccessKeyID := "AKIABLABLA"
		expectedSecretAccessKey := "SMAMSLscSercreasdas"
		expectedRegion := "eu-west-1"

		params := map[string]string{
			"accessKeyID":     expectedAccessKeyID,
			"secretAccessKey": expectedSecretAccessKey,
			"region":          expectedRegion,
		}

		Convey("When creating AWS Config from them", func() {
			config, err := awsConfigFromParams(params)
			So(err, ShouldBeNil)
			Convey("Credentials are created correctly", func() {
				actualCredentials, err := config.Credentials.Get()
				So(err, ShouldBeNil)
				So(aws.StringValue(config.Region), ShouldEqual, expectedRegion)
				So(actualCredentials.AccessKeyID, ShouldEqual, expectedAccessKeyID)
				So(actualCredentials.SecretAccessKey, ShouldEqual, expectedSecretAccessKey)
			})
		})
	})
}
