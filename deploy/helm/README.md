# External Secret Operator Helm Chart

This chart packages the [External Secret Operator](https://github.com/ContainerSolutions/externalsecret-operator) which makes easy to inject information stored in password managers such as [AWS Secret Manager](https://aws.amazon.com/secrets-manager/) in your cluster as Kubernetes [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

## Prerequisites

* Kubernetes 1.7+ with Custom Resource Definition support

## Chart Details

This chart will create:

1. A [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) with a single replica running the External Secret Operator

1. A [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) named `ExternalSecret` that will be used to create `Secret` resources fetching information from the password manager
 
3. A [Secret](https://kubernetes.io/docs/concepts/configuration/secret/) (optional) that will hold the credentials the External Secret Operator will use to authenticate with an external password manager

4. A [ServiceAccount](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/) (optional) which will be used by the External Secret Operator

5. Two [Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole)/[RoleBinding](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding) pairs (optional), one to restrict the actions the External Secret Operator caniperform in the namespace in which is running and the other to restrict the operations it can perform in the namespace which is watching

### Installation

The following command will install the Exernal Secret Operator using the AWS Secrets Manager (asm) backend:

```shell
helm upgrade --install asm1 --wait \
    --set operatorName="asm-example" \
    --set secret.data.Type="asm" \
    --set secret.data.Parameters.accessKeyID="$AWS_ACCESS_KEY_ID" \
    --set secret.data.Parameters.region="$AWS_DEFAULT_REGION" \
    --set secret.data.Parameters.secretAccessKey="$AWS_SECRET_ACCESS_KEY" \
    ./deploy/helm/.
```

### Configuration

You can provide the configuration of the External Secret Operator using the `secret` key in `values.yaml`. 

|Parameter|Description|Default|
| - | - | - |
| `replicaCount` | Number of replicas to run | `1`
| `image.repository` | Image repository | `containersol/externalsecret-operator`
| `image.tag` | Image tag | `0.2.0`
| `image.pullPolicy` | Image pull policy | `IfNotPresent`
| `watchNamespace` | Namespace to watch for `ExternalSecret` resources. If empty, will be the same as the one where the operator will be deployed | `""`
| `operatorName` | Name passed as `OPERATOR_NAME` environment variable. Referenced by `ExternalSecret` resources in `Backend` field | `externalsecret-operator`
| `secret.create` | Whether the secret containing the operator configuration should be created | `true`
| `secret.name` | Name of the Secret that contains the operator configuration. If empty and `secret.create` is `true`, a secret based on the release name will be generated | `""`
| `secret.key` | Key in the secret holding the operator configuration | `config.json`
| `secret.data` | External Secret Operator configuration. This depend on the backend type | [External Secret Operator Default Configuration](#markdown-header-default-pubsub-configuration)

#### External Secret Operator Default Configuration

The following default configuration is added in the secret if no other is specified:

```yaml
data:
  Type: "dummy"
  Parameters:
    suffix: "-externalsecretsoperatorwashere"
```

This will configure the operator to use the `dummy` backend which simply adds a suffix to the `ExternalSecret` key field and uses that as a value for the generated `Secret`. It is normally used for testing purposes.
