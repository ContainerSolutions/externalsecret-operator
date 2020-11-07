// Package gitlab implements a gitlab backend that can be used to inject
// Gitlab CI/CD variables as kubernetes secrets.
package gitlab

import (
	"encoding/json"
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
	client    *gitlab.Client
	projectID interface{}
}

func init() {
	backend.Register("gitlab", NewBackend)
}

// NewBackend gives you a new gitlab backend
func NewBackend() backend.Backend {
	return &Backend{}
}

// Init implements SecretsBackend interface
func (d *Backend) Init(parameters map[string]interface{}, credentials []byte) error {
	var err error

	if len(parameters) == 0 {
		log.Error(fmt.Errorf("error"), "empty or invalid parameters: ")
		return fmt.Errorf("empty or invalid parameters")
	}

	gitlabCreds := &GitlabCredentials{}
	if err := json.Unmarshal(credentials, gitlabCreds); err != nil {
		log.Error(err, "Unmarshalling failed")
		return &ErrInitFailed{message: err.Error()}
	}

	baseURL, ok := parameters["baseURL"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing baseURL parameter: ")
		return fmt.Errorf("missing baseURL parameter")
	}

	projectID, ok := parameters["projectID"]
	if !ok {
		log.Error(fmt.Errorf("error"), "missing projectID parameter: ")
		return fmt.Errorf("missing projectID parameter")
	}

	d.projectID = projectID
	d.client, err = gitlab.NewClient(gitlabCreds.Token, gitlab.WithBaseURL(baseURL.(string)))
	if err != nil {
		log.Error(fmt.Errorf("error"), "failed to create client: ")
		return fmt.Errorf("failed to create client")
	}
	return nil
}

// Get takes a key and version, and returns the value
func (d *Backend) Get(key string, version string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key provided")
	}

	variable, _, err := d.client.ProjectVariables.GetVariable(fmt.Sprintf("%.f", d.projectID), key, nil)
	if err != nil {
		return "", err
	}

	log.Info("Get was successful for the Gitlab")

	return variable.Value, nil
}

type GitlabCredentials struct {
	Token string `json:"token"`
}
