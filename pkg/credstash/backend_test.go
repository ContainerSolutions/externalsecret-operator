package credstash

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	. "github.com/smartystreets/goconvey/convey"
	unicreds "github.com/versent/unicreds"
)

type mockedSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
	withError bool
}

// GetHighestVersionSecret mocked to return expected value
func (s mockedSecretsManager) GetHighestVersionSecret(tableName *string, name string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error) {
	var cred *unicreds.Credential
	return &unicreds.DecryptedCredential{Credential: cred, Secret: "secretValue"}, nil
}

// GetSecret mocked to return expected value
func (s mockedSecretsManager) GetSecret(tableName *string, name string, version string, encContext *unicreds.EncryptionContextValue) (*unicreds.DecryptedCredential, error) {
	var cred *unicreds.Credential
	return &unicreds.DecryptedCredential{Credential: cred, Secret: "secretValue"}, nil
}

// SetKMSConfig sets configuration for KMS access
func (s mockedSecretsManager) SetKMSConfig(config *aws.Config) {
}

// SetDynamoDBConfig sets configuration for DynamoDB access
func (s mockedSecretsManager) SetDynamoDBConfig(config *aws.Config) {
}

func TestNewBackend(t *testing.T) {
	Convey("When creating a new Credstash backend", t, func() {
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

	Convey("Given an uninitialized CredstashSecretsManagerBackend", t, func() {
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
		backend.SecretsManager = mockedSecretsManager{}
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}

type credentialsAndParametersTest struct {
	credentials               string
	parameters                map[string]interface{}
	expectedAccessKeyID       string
	expectedRegion            string
	expectedTable             string
	expectedEncryptionContext map[string]string
	expectedSecretAccessKey   string
	expectedSessionToken      string
	expectedErrorAssertion    func(interface{}, ...interface{}) string
	expectedErrorString       string
}

func TestInit(t *testing.T) {

	tests := []credentialsAndParametersTest{
		{
			credentials: `{
				"accessKeyID":     "AKIABLABLA",
				"secretAccessKey": "CredSsecrets",
				"sessionToken": ""
			}`,
			parameters: map[string]interface{}{
				"region": "eu-mediterranean-1",
				"table":  "credential-store",
				"encryptionContext": map[string]string{
					"securityKey": "securityValue",
				},
			},
			expectedAccessKeyID: "AKIABLABLA",
			expectedRegion:      "eu-mediterranean-1",
			expectedTable:       "credential-store",
			expectedEncryptionContext: map[string]string{
				"securityKey": "securityValue",
			},
			expectedSecretAccessKey: "CredSsecrets",
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
				"table":  "credential-store",
				"encryptionContext": map[string]string{
					"securityKey": "securityValue",
				},
			},
			expectedAccessKeyID: "AKIABLABLA",
			expectedRegion:      "eu-mediterranean-1",
			expectedTable:       "credential-store",
			expectedEncryptionContext: map[string]string{
				"securityKey": "securityValue",
			},
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
				"table":  "credential-store",
				"encryptionContext": map[string]string{
					"securityKey": "securityValue",
				},
			},
			expectedAccessKeyID: "some",
			expectedRegion:      "other",
			expectedTable:       "credential-store",
			expectedEncryptionContext: map[string]string{
				"securityKey": "securityValue",
			},
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
				"table":  "credential-store",
				"encryptionContext": map[string]string{
					"securityKey": "securityValue",
				},
			},
			expectedAccessKeyID: "some",
			expectedRegion:      "eu-west-2",
			expectedTable:       "credential-store",
			expectedEncryptionContext: map[string]string{
				"securityKey": "securityValue",
			},
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
				"table":  "credential-store",
				"encryptionContext": map[string]string{
					"securityKey": "securityValue",
				},
			},
			expectedAccessKeyID: "some",
			expectedRegion:      "eu-west-2",
			expectedTable:       "credential-store",
			expectedEncryptionContext: map[string]string{
				"securityKey": "securityValue",
			},
			expectedSecretAccessKey: "VGhhdEtleQoVGhhdEtleQo",
			expectedSessionToken:    "tZWtletZWtle",
		},
	}

	for _, test := range tests {
		Convey("Given AWS credentials", t, func() {

			Convey("When initializing an Credstash backend", func() {
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
					"secretAccessKey": "CredSsecrets",
					"sessionToken": ""
				}`,
			parameters:              map[string]interface{}{},
			expectedAccessKeyID:     "AKIABLABLA",
			expectedRegion:          "eu-mediterranean-1",
			expectedSecretAccessKey: "CredSsecrets",
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
			expectedSecretAccessKey: "CredSsecrets",
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
