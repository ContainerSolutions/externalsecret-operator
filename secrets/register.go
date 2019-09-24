package secrets

import (
	// Register your backends here
	_ "github.com/containersolutions/externalsecretoperator/secrets/asm"
	_ "github.com/containersolutions/externalsecretoperator/secrets/dummy"
	_ "github.com/containersolutions/externalsecretoperator/secrets/onepassword"
)
