// Package secrets implements the logic and data structures to handle external
// secrets backends. Each external secret implementation will reside in its own
// package. A "Dummy" backend is provided as reference.
//
// Backends must register their "type" using a function to instantiate themselves. An
// easy way of doing so is calling the secrets.Register function inside the
// backend package init() function:
//
//		func init() {
//			secrets.Register("dummy", NewBackend)
//		}
//
//		// NewBackend gives you an new Dummy Backend
//		func NewBackend() secrets.Backend {
//			return &Backend{}
//		}
//
package secrets
