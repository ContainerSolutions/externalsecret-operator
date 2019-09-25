package secrets

import (
	// Register your backends here
	_ "github.com/containersolutions/externalsecret-operator/secrets/asm"
	_ "github.com/containersolutions/externalsecret-operator/secrets/dummy"
	_ "github.com/containersolutions/externalsecret-operator/secrets/onepassword"
)
