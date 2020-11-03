package backend

import (
	"encoding/json"
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("config")

//ConfigEnvVar holds the name of the Environment Variable scanned for config
const ConfigEnvVar string = "OPERATOR_CONFIG"

//Config represent configuration information for the secrets backend
type Config struct {
	Type       string
	Parameters map[string]interface{}
	Auth       map[string]interface{}
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

// ConfigFromCtrl returns a Config object based on the byte data passed as parameter
func ConfigFromCtrl(data []byte) (*Config, error) {
	backendConfig := &Config{}
	err := json.Unmarshal(data, backendConfig)
	if err != nil {
		return nil, err
	}
	return backendConfig, nil
}

//ConfigFromEnv parses Config from environment variable
func ConfigFromEnv() (*Config, error) {
	data, present := os.LookupEnv(ConfigEnvVar)
	if !present {
		return nil, fmt.Errorf("cannot find config: `%v` not set", ConfigEnvVar)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("cannot find config: `%v` not set", ConfigEnvVar)
	}
	return ConfigFromJSON(data)
}
