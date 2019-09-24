// build +integration
package onepassword

import (
	"fmt"
	"strings"
	"testing"

	"os"

	ioutil "io/ioutil"

	"github.com/ContainerSolutions/externalsecret-operator/secrets/backend"
)

func TestOnePasswordBackend(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Given
	_ = backend.InitFromEnv()
	backend.Register("onepassword", NewBackend)
	backend := backend.Instances["onepassword"]

	secretKey, expectedValue := GetKeyAndValue(t)

	// When
	value, _ := backend.Get(secretKey)

	// Then
	if expectedValue != value {
		fmt.Printf("Expected value '%s' is not equal to value '%s'\n", expectedValue, value)
		t.Fail()
	}
}

func GetKeyAndValue(t *testing.T) (string, string) {
	envVar := "SECRET_KEY_FILE"
	secretKeyFile := CheckAndGetenv(envVar, t)

	if secretKeyFile == "" || !strings.HasPrefix(secretKeyFile, "secret-") {
		fmt.Printf("env var '%s' should point to file whose filename consists of 'secret-' plus the key of the secret. The contents of the file should be the secret value.\n", envVar)
		t.Fail()
	}

	bytes, _ := ioutil.ReadFile(secretKeyFile)
	value := string(bytes)

	key := strings.TrimPrefix(secretKeyFile, "secret-")

	return key, value
}

func CheckAndGetenv(name string, t *testing.T) string {
	value := os.Getenv(name)
	if value == "" {
		fmt.Printf("please specify '%s' env var\n", name)
		t.Fail()
	}
	return value
}
