package secrets

//Backends is a map of backends
var Backends map[string]BackendIface

//Backend is a Backend to a Secret Store
type Backend struct {
	Name string
}

//BackendIface is an interface to a Backend
type BackendIface interface {
	Init(...interface{}) error
	Get(string) (string, error)
}

//BackendRegister adds a Backend to the Backends map
func BackendRegister(name string, backend BackendIface) {
	if Backends == nil {
		Backends = make(map[string]BackendIface)
	}

	Backends[name] = backend
}
