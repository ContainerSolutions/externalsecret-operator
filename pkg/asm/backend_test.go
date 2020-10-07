package asm

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	. "github.com/smartystreets/goconvey/convey"
)

type mockedSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
	withError bool
}

func (m *mockedSecretsManager) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	mockSecretString := *input.SecretId + "Value"
	output := &secretsmanager.GetSecretValueOutput{
		Name:         input.SecretId,
		SecretString: &mockSecretString,
	}
	if m.withError {
		return output, errors.New("oops")
	}
	return output, nil
}

func TestNewBackend(t *testing.T) {
	Convey("When creating a new ASM backend", t, func() {
		backend := NewBackend()
		So(backend, ShouldNotBeNil)
		So(backend, ShouldHaveSameTypeAs, &Backend{})
	})
}

func TestGet(t *testing.T) {
	secretKey := "secret"
	keyVersion := ""
	secretValue := "secretValue"
	expectedValue := secretValue

	Convey("Given an uninitialized AWSSecretsManagerBackend", t, func() {
		backend := Backend{}
		Convey("When retrieving a secret", func() {
			_, err := backend.Get(secretKey, keyVersion)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "backend not initialized")
			})
		})
	})

	Convey("Given an initialized AWSSecretsManagerBackend", t, func() {
		backend := Backend{}
		backend.SecretsManager = &mockedSecretsManager{}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})

	Convey("Given an initialized AWSSecretsManagerBackend (withError: true)", t, func() {
		backend := Backend{}
		backend.SecretsManager = &mockedSecretsManager{withError: true}
		Convey("When retrieving a secret", func() {
			_, err := backend.Get(secretKey, keyVersion)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

type envVariablesTest struct {
	envVariables            map[string]string
	parameters              map[string]string
	expectedAccessKeyID     string
	expectedRegion          string
	expectedSecretAccessKey string
	expectedErrorAssertion  func(interface{}, ...interface{}) string
	expectedErrorString     string
}

func TestInit(t *testing.T) {
	// https://docs.aws.amazon.com/sdk-for-go/api/aws/session/

	tests := []envVariablesTest{
		{
			envVariables: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIABLABLA",
				"AWS_REGION":            "eu-mediterranean-1",
				"AWS_SECRET_ACCESS_KEY": "SMMSsecrets",
			},
			parameters:              nil,
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "SMMSsecrets",
		},
		{
			envVariables: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIABLABLA",
				"AWS_REGION":            "eu-mediterranean-1",
				"AWS_SECRET_ACCESS_KEY": "SMMSsecrets",
			},
			parameters:              map[string]string{},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "SMMSsecrets",
		},
		{
			envVariables: map[string]string{
				"AWS_ACCESS_KEY_ID":     "AKIABLABLA",
				"AWS_SECRET_ACCESS_KEY": "eu-mediterranean-1",
				"AWS_REGION":            "SMMSsecrets",
			},
			parameters: map[string]string{
				"accessKeyID":     "some",
				"region":          "other",
				"secretAccessKey": "value",
			},
			expectedAccessKeyID:     "some",
			expectedRegion:          "other",
			expectedSecretAccessKey: "value",
		},
	}

	for _, test := range tests {
		Convey("Given AWS credentials using environment variables", t, func() {
			for k, v := range test.envVariables {
				os.Setenv(k, v)
			}
			Convey("When initializing an ASM backend", func() {
				b := Backend{}
				err := b.Init(test.parameters)
				So(err, ShouldBeNil)
				Convey("Then credentials are reflected in the AWS session", func() {
					actualCredentials, err := b.session.Config.Credentials.Get()
					So(err, ShouldBeNil)
					So(*b.session.Config.Region, ShouldEqual, test.expectedRegion)
					So(actualCredentials.AccessKeyID, ShouldEqual, test.expectedAccessKeyID)
					So(actualCredentials.SecretAccessKey, ShouldEqual, test.expectedSecretAccessKey)
				})
			})

		})
	}
}
