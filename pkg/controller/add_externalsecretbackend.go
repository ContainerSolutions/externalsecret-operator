package controller

import (
	"github.com/ContainerSolutions/externalsecret-operator/pkg/controller/externalsecretbackend"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, externalsecretbackend.Add)
}
