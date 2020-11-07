// Package dummy implements an example backend that can be used for testing
// purposes. It acceps a "suffix" as a parameter that will be appended to the
// key passed to the Get function.
package gitlab

import (
	"fmt"

	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	gitlab "github.com/xanzy/go-gitlab"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("gitlab")

type ErrInitFailed struct {
	message string
}

func (e *ErrInitFailed) Error() string {
	return fmt.Sprintf("gitlab backend init failed: %s", e.message)
}

type ErrGet struct {
	itemName string
	message  string
}

func (e *ErrGet) Error() string {
	return fmt.Sprintf("gitlab backend get '%s' failed: %s", e.itemName, e.message)
}

// Backend is a gitlab variables backend
type Backend struct {
	client *gitlab.Client
	pid    interface{}
	suffix string
}

func init() {
	backend.Register("gitlab", NewBackend)
}

// NewBackend gives you a new gitlab backend
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init implements SecretsBackend interface, sets the suffix
func (d *Backend) Init(parameters map[string]interface{}, credentials []byte) error {
	var err error

	if len(parameters) == 0 {
		log.Error(fmt.Errorf("error"), "empty or invalid parameters: ")
		return fmt.Errorf("empty or invalid parameters")
	}

	suffix, ok := parameters["Suffix"].(string)
	if !ok {
		log.Error(fmt.Errorf("error"), "missing parameters: ")
		return fmt.Errorf("missing parameters")
	}

	d.suffix = suffix

	// initialize gitalb client
	token, ok := parameters["token"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing token parameter: ")
		return fmt.Errorf("missing token parameter")
	}

	baseURL, ok := parameters["baseURL"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing baseURL parameter: ")
		return fmt.Errorf("missing baseURL parameter")
	}

	pid, ok := parameters["pid"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing pid parameter: ")
		return fmt.Errorf("missing pid parameter")
	}

	d.pid = pid
	d.client, err = gitlab.NewClient(token.(string), gitlab.WithBaseURL(baseURL.(string)))
	if err != nil {
		log.Error(fmt.Errorf("error"), "failed to create client: ")
		return fmt.Errorf("failed to create client")
	}
	return nil
}

// Get takes a key and version, and returns the value
func (d *Backend) Get(key string, version string) (string, error) {
	if d.suffix == "" {
		return "", fmt.Errorf("backend is not initialized")
	}

	if key == "" {
		return "", fmt.Errorf("empty key provided")
	}

	variable, _, err := d.client.ProjectVariables.GetVariable(d.pid, key, nil)
	if err != nil {
		return "", err
	}
	return variable.Value, nil
}
