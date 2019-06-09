package secrets

import (
	"fmt"
	"sync"
)

// Backend is an abstract backend interface
type Backend interface {
	Init(map[string]string) error
	Get(string) (string, error)
}

// Instances are instantiated secret backends
var Instances map[string]Backend

// Functions is a map of labelled functions that return secret backend instances
var Functions map[string]func() Backend

var initLock sync.Mutex

// Instantiate instantiates a Backend of type `backendType`
func Instantiate(name string, backendType string) error {
	if Instances == nil {
		Instances = make(map[string]Backend)
	}

	function, found := Functions[backendType]
	if !found {
		return fmt.Errorf("unknown backend type: '%v'", backendType)
	}

	Instances[name] = function()

	return nil
}

// Register registers a new backend type with name `name`
// function is a function that returns a backend of that type
func Register(name string, function func() Backend) {
	if Functions == nil {
		Functions = make(map[string]func() Backend)
	}

	Functions[name] = function
}

// InitFromEnv initializes a backend looking into Env for config data
func InitFromEnv() error {
	initLock.Lock()
	defer initLock.Unlock()

	config, err := ConfigFromEnv()
	if err != nil {
		return err
	}

	err = Instantiate(config.Name, config.Type)
	if err != nil {
		return err
	}

	err = Instances[config.Name].Init(config.Parameters)

	return err
}
