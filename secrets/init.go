package secrets

import (
	// Register your backends here
	_ "github.com/ContainerSolutions/externalsecret-operator/secrets/asm"
	_ "github.com/ContainerSolutions/externalsecret-operator/secrets/dummy"
	_ "github.com/ContainerSolutions/externalsecret-operator/secrets/onepassword"
)
