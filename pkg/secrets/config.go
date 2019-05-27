package secrets

import (
	"encoding/json"
	"fmt"
	"os"
)

//ConfigEnvVar holds the name of the Environment Variable scanned for config
const ConfigEnvVar string = "BACKEND_CONFIG"

//BackendConfig represent configuration information for secret backend
type BackendConfig struct {
	Name       string
	Type       string
	Parameters map[string]string
}

// BackendConfigFromJSON returns a BackendConfig object based on the string data passed as parameter
func BackendConfigFromJSON(data string) (*BackendConfig, error) {
	backendConfig := &BackendConfig{}
	err := json.Unmarshal([]byte(data), backendConfig)
	if err != nil {
		return nil, err
	}
	return backendConfig, nil
}

//BackendConfigFromEnv parse BackendConfiguration from environment variable
func BackendConfigFromEnv() (*BackendConfig, error) {
	data := os.Getenv(ConfigEnvVar)
	if len(data) == 0 {
		return nil, fmt.Errorf("Cannot find config: `%v` not set", ConfigEnvVar)
	}
	return BackendConfigFromJSON(data)
}
