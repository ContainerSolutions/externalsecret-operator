package backend

import (
	"fmt"
	"strings"
	"sync"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("backend")

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

	log.Info("instantiate", "name", name, "type", backendType)
	Instances[name] = function()

	return nil
}

// Register registers a new backend type with name `name`
// function is a function that returns a backend of that type
func Register(name string, function func() Backend) {
	if Functions == nil {
		Functions = make(map[string]func() Backend)
	}

	log.Info("register", "type", name)
	Functions[name] = function
}

// InitFromEnv initializes a backend looking into Env for config data
func InitFromEnv() error {
	initLock.Lock()
	defer initLock.Unlock()
	log.Info("initFromEnv", "availableBackends", strings.Join(availableBackends(), ","))

	config, err := ConfigFromEnv()
	if err != nil {
		return err
	}

	operatorName, err := k8sutil.GetOperatorName()
	if err != nil {
		return err
	}

	err = Instantiate(operatorName, config.Type)
	if err != nil {
		return err
	}

	log.Info("initialize", "name", operatorName)
	err = Instances[operatorName].Init(config.Parameters)

	return err
}

func availableBackends() []string {
	backends := []string{}
	for k := range Functions {
		backends = append(backends, k)
	}
	return backends
}
