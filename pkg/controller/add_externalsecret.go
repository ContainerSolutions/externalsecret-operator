package controller

import (
	"github.com/ContainerSolutions/externalsecret-operator/pkg/controller/externalsecret"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, externalsecret.Add)
}
