package controller

import (
	// Register your backends here
	_ "github.com/containersolutions/externalsecret-operator/pkg/asm"
	_ "github.com/containersolutions/externalsecret-operator/pkg/dummy"
	_ "github.com/containersolutions/externalsecret-operator/pkg/onepassword"
)
