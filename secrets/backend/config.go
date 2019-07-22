package backend

import (
	"encoding/json"
	"fmt"
	"os"
)

//ConfigEnvVar holds the name of the Environment Variable scanned for config
const ConfigEnvVar string = "OPERATOR_CONFIG"

//Config represent configuration information for the secrets backend
type Config struct {
	Type       string
	Parameters map[string]string
}

// ConfigFromJSON returns a Config object based on the string data passed as parameter
func ConfigFromJSON(data string) (*Config, error) {
	backendConfig := &Config{}
	err := json.Unmarshal([]byte(data), backendConfig)
	if err != nil {
		return nil, err
	}
	return backendConfig, nil
}

//ConfigFromEnv parses Config from environment variable
func ConfigFromEnv() (*Config, error) {
	data := os.Getenv(ConfigEnvVar)
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot find config: `%v` not set", ConfigEnvVar)
	}
	return ConfigFromJSON(data)
}
