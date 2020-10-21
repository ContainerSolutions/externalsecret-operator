package backend

import (
	"fmt"
	"strings"
	"sync"

	config "github.com/containersolutions/externalsecret-operator/pkg/config"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("backend")

// Backend is an abstract backend interface
type Backend interface {
	Init(map[string]interface{}, []byte) error
	Get(string, string) (string, error)
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
		log.Error(fmt.Errorf("error"), fmt.Sprintf("unknown backend type: '%v'", backendType))
		return fmt.Errorf("unknown backend type: '%v'", backendType)
	}

	log.Info("Instantiate", "name", name, "type", backendType)
	Instances[name] = function()

	return nil
}

// Register registers a new backend type with name `name`staging
// function is a function that returns a backend of that type
func Register(name string, function func() Backend) {
	if Functions == nil {
		Functions = make(map[string]func() Backend)
	}

	log.Info("Register", "type", name)
	Functions[name] = function
}

// InitFromEnv initializes a backend looking into Env for config data
func InitFromEnv(leaderID string) error {
	initLock.Lock()
	defer initLock.Unlock()
	log.Info("InitFromEnv", "availableBackends", strings.Join(availableBackends(), ","))

	config, err := config.ConfigFromEnv()
	if err != nil {
		return err
	}

	err = Instantiate(leaderID, config.Type)
	if err != nil {
		log.Error(err, "")
		return err
	}

	log.Info("Initialize", "name", leaderID)
	err = Instances[leaderID].Init(config.Parameters, []byte(""))

	return err
}

// InitFromCtrl initializes within a controller
func InitFromCtrl(contrl string, config *config.Config, credentials []byte) error {
	initLock.Lock()
	defer initLock.Unlock()
	log.Info("InitFromCtrl", "availableBackends", strings.Join(availableBackends(), ","))

	err := Instantiate(contrl, config.Type)
	if err != nil {
		log.Error(err, "")
		return err
	}

	log.Info("Initialize", "name", contrl)
	err = Instances[contrl].Init(config.Parameters, credentials)

	return err
}

func availableBackends() []string {
	backends := []string{}
	for k := range Functions {
		backends = append(backends, k)
	}
	return backends
}
