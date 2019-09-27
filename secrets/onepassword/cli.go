package onepassword

import (
	op "github.com/ameier38/onepassword"
)

type Cli interface {
	GetItem(op.VaultName, op.ItemName) (op.ItemMap, error)
}
