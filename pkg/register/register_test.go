package controller

import (
	"testing"

	"github.com/containersolutions/externalsecret-operator/pkg/backend"
)

var expectedRegisteredBackends = []string{
	"asm",
	"dummy",
	"onepassword",
	"gsm",
}

func TestInit(t *testing.T) {
	for _, k := range expectedRegisteredBackends {
		_, found := backend.Functions[k]
		if !found {
			t.Errorf("registered backend expected but not found: '%v'", k)
		}
	}
}
