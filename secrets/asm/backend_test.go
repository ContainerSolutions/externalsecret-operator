package asm

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
		backend := Backend{}
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

type parametersTest struct {
	parameters              map[string]string
	expectedAccessKeyID     string
	expectedRegion          string
	expectedSecretAccessKey string
	expectedErrorAssertion  func(interface{}, ...interface{}) string
	expectedErrorString     string
}

func TestAWSConfigFromParams(t *testing.T) {
	expectedAccessKeyID := "AKIABLABLA"
	expectedSecretAccessKey := "SMAMSLscSercreasdas"
	expectedRegion := "eu-west-1"

	tests := []parametersTest{
		{
			parameters: map[string]string{
				"accessKeyID":     expectedAccessKeyID,
				"region":          expectedRegion,
				"secretAccessKey": expectedSecretAccessKey,
			},
			expectedAccessKeyID:     expectedAccessKeyID,
			expectedRegion:          expectedRegion,
			expectedSecretAccessKey: expectedSecretAccessKey,
			expectedErrorAssertion:  ShouldBeNil,
		},

		{
			parameters: map[string]string{
				"accessKeyID":     expectedAccessKeyID,
				"secretAccessKey": expectedSecretAccessKey,
			},
			expectedAccessKeyID:     expectedAccessKeyID,
			expectedRegion:          expectedRegion,
			expectedSecretAccessKey: expectedSecretAccessKey,
			expectedErrorAssertion:  ShouldNotBeNil,
			expectedErrorString:     "Invalid init parameters: expected `region` not found",
		},
		{
			parameters: map[string]string{
				"region":          expectedRegion,
				"secretAccessKey": expectedSecretAccessKey,
			},
			expectedAccessKeyID:     expectedAccessKeyID,
			expectedRegion:          expectedRegion,
			expectedSecretAccessKey: expectedSecretAccessKey,
			expectedErrorAssertion:  ShouldNotBeNil,
			expectedErrorString:     "Invalid init parameters: expected `accessKeyID` not found",
		},

		{
			parameters: map[string]string{
				"accessKeyID": expectedAccessKeyID,
				"region":      expectedRegion,
			},
			expectedAccessKeyID:     expectedAccessKeyID,
			expectedRegion:          expectedRegion,
			expectedSecretAccessKey: expectedSecretAccessKey,
			expectedErrorAssertion:  ShouldNotBeNil,
			expectedErrorString:     "Invalid init parameters: expected `secretAccessKey` not found",
		},
	}

	for _, test := range tests {
		Convey("Given a set of params", t, func() {
			Convey("When creating AWS Config from them", func() {
				config, err := awsConfigFromParams(test.parameters)
				So(err, test.expectedErrorAssertion)
				if err != nil {
					So(err.Error(), ShouldEqual, test.expectedErrorString)
				} else {
					Convey("Credentials are created correctly", func() {
						actualCredentials, err := config.Credentials.Get()
						So(err, ShouldBeNil)
						So(aws.StringValue(config.Region), ShouldEqual, test.expectedRegion)
						So(actualCredentials.AccessKeyID, ShouldEqual, test.expectedAccessKeyID)
						So(actualCredentials.SecretAccessKey, ShouldEqual, test.expectedSecretAccessKey)
					})
				}
			})
		})
	}

}
