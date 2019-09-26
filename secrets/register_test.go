package secrets

import (
	"testing"

	"github.com/containersolutions/externalsecret-operator/secrets/backend"
)

var expectedRegisteredBackends = []string{
	"asm",
	"dummy",
	"onepassword",
}

func TestInit(t *testing.T) {
	for _, k := range expectedRegisteredBackends {
		_, found := backend.Functions[k]
		if !found {
			t.Errorf("registered backend expected but not found: '%v'", k)
		}
	}
}
