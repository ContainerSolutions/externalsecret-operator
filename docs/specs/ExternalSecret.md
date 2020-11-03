```
apiVersion: secrets.externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecret
metadata: {...}
spec:
  # Optional
  # The amount of time before the values will be read again from the store
  # Secret Rotation Period;
  # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  refreshInterval: [String] - Default value "1h"

  # Secret name to be created by ExternalSecret
  # Optional
  target: 
    # The secret name of the resource
    # defaults to .metadata.name of the ExternalSecret. immutable.
    name: my-secret

  # Required 
  # A reference to the store used to fetch the secrets
  storeRef:
    kind: SecretStore # ClusterSecretStore
    name: my-store

  # Required
  # data contains key/value pairs which correspond to the keys in the resulting secret
  data: [Array]
    - key: [String]
      version: [String]
    
status: {}
```