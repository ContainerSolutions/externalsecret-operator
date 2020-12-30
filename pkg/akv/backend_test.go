package akv

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
)

type mockedClient struct {
}

func TestNewBackend(t *testing.T) {
	backend := Backend{}
	_, err := backend.Get("hello", "")

	if err == nil {
		t.Errorf("There should have been an error because the backend has not been initialized")
		return
	}

	if err.Error() != "Azure Key Vault backend not initialized" {
		t.Error(err)
		return
	}
}

var flagtests = []struct {
	in   string
	out  string
	pass bool
}{
	{"hello", "hello", true},
	{"hello", "world", false},
	{"world", "world", true},
	{"foo", "bar", false},
}

func TestGetSecret(t *testing.T) {

	backend := Backend{}
	backend.Client = &mockedClient{}

	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {

			result, err := backend.Get(tt.in, "")

			if err != nil {
				t.Error(err)
			} else if (result == tt.out) != tt.pass {
				t.Errorf("Expected: %s, got: %s", tt.out, result)
			}
		})
	}
}

func (m *mockedClient) GetSecret(context context.Context, url string, key string, version string) (keyvault.SecretBundle, error) {

	return keyvault.SecretBundle{Value: &key}, nil
}
