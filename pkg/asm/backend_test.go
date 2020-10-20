package asm

import (
	"errors"
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
	mockedSecretBinary := []byte("b2ggbm8gVGhleSBjYW4gc2VlIHVzIG5vdw==")

	output := &secretsmanager.GetSecretValueOutput{
		Name: input.SecretId,
	}

	if *input.SecretId == "secretKeyBinary" {
		output.SecretBinary = mockedSecretBinary
	} else {
		output.SecretString = &mockSecretString
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
	secretKeyBinary := "secretKeyBinary"
	keyVersion := ""
	secretValue := "secretValue"
	expectedValue := secretValue
	expectedSecretBinaryValue := "oh no They can see us now"

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

	Convey("Given an initialized AWSSecretsManagerBackend", t, func() {
		backend := Backend{}
		backend.SecretsManager = &mockedSecretsManager{}
		Convey("When retrieving a binary secret", func() {
			actualValue, err := backend.Get(secretKeyBinary, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedSecretBinaryValue)
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

type credentialsAndParametersTest struct {
	credentials             string
	parameters              map[string]interface{}
	expectedAccessKeyID     string
	expectedRegion          string
	expectedSecretAccessKey string
	expectedSessionToken    string
	expectedErrorAssertion  func(interface{}, ...interface{}) string
	expectedErrorString     string
}

func TestInit(t *testing.T) {
	// https://docs.aws.amazon.com/sdk-for-go/api/aws/session/

	tests := []credentialsAndParametersTest{
		{
			credentials: `{
				"accessKeyID":     "AKIABLABLA",
				"secretAccessKey": "SMMSsecrets",
				"sessionToken": ""
			}`,
			parameters: map[string]interface{}{
				"region": "eu-mediterranean-1",
			},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "SMMSsecrets",
			expectedSessionToken:    "",
		},
		{
			credentials: `{
				"accessKeyID":     "AKIABLABLA",
				"secretAccessKey": "QW5vdGhlcmtleQoQW5vdGhlcmtleQo",
				"sessionToken": ""
			}`,
			parameters: map[string]interface{}{
				"region": "eu-mediterranean-1",
			},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "QW5vdGhlcmtleQoQW5vdGhlcmtleQo",
			expectedSessionToken:    "",
		},
		{
			credentials: `{
				"accessKeyID":     "some",
				"secretAccessKey": "U29tZWtleQoU29tZWtleQo",
				"sessionToken": ""
			}`,
			parameters: map[string]interface{}{
				"region": "other",
			},
			expectedAccessKeyID:     "some",
			expectedRegion:          "other",
			expectedSecretAccessKey: "U29tZWtleQoU29tZWtleQo",
			expectedSessionToken:    "",
		},

		{
			credentials: `{
				"accessKeyID":     "some",
				"secretAccessKey": "VGhhdEtleQoVGhhdEtleQo",
				"sessionToken": "EtleQoVGhhEtleQoVGhh"
			}`,
			parameters: map[string]interface{}{
				"region": "eu-west-2",
			},
			expectedAccessKeyID:     "some",
			expectedRegion:          "eu-west-2",
			expectedSecretAccessKey: "VGhhdEtleQoVGhhdEtleQo",
			expectedSessionToken:    "EtleQoVGhhEtleQoVGhh",
		},

		{
			credentials: `{
				"accessKeyID":     "some",
				"secretAccessKey": "VGhhdEtleQoVGhhdEtleQo",
				"sessionToken": "tZWtletZWtle"
			}`,
			parameters: map[string]interface{}{
				"region": "",
			},
			expectedAccessKeyID:     "some",
			expectedRegion:          "eu-west-2",
			expectedSecretAccessKey: "VGhhdEtleQoVGhhdEtleQo",
			expectedSessionToken:    "tZWtletZWtle",
		},
	}

	for _, test := range tests {
		Convey("Given AWS credentials", t, func() {

			Convey("When initializing an ASM backend", func() {
				b := Backend{}
				err := b.Init(test.parameters, []byte(test.credentials))
				So(err, ShouldBeNil)
				Convey("Then credentials are reflected in the AWS session", func() {
					actualCredentials, err := b.session.Config.Credentials.Get()
					So(err, ShouldBeNil)
					So(*b.session.Config.Region, ShouldEqual, test.expectedRegion)
					So(actualCredentials.AccessKeyID, ShouldEqual, test.expectedAccessKeyID)
					So(actualCredentials.SecretAccessKey, ShouldEqual, test.expectedSecretAccessKey)
					So(actualCredentials.SessionToken, ShouldEqual, test.expectedSessionToken)
				})
			})

		})
	}

	Convey("When missing region parameter", t, func() {
		testParams := credentialsAndParametersTest{
			credentials: `{
					"accessKeyID":     "AKIABLABLA",
					"secretAccessKey": "SMMSsecrets",
					"sessionToken": ""
				}`,
			parameters:              map[string]interface{}{},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "SMMSsecrets",
			expectedSessionToken:    "",
		}

		b := Backend{}
		err := b.Init(testParams.parameters, []byte(testParams.credentials))
		Convey("Then an error is returned", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "AWS region parameter missing")
		})
	})

	Convey("When invalid credentials are passed", t, func() {
		testParams := credentialsAndParametersTest{
			credentials:             "",
			parameters:              map[string]interface{}{},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "SMMSsecrets",
			expectedSessionToken:    "",
		}

		b := Backend{}
		err := b.Init(testParams.parameters, []byte(testParams.credentials))
		Convey("Then an error is returned", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unexpected end of JSON input")
		})
	})
}
