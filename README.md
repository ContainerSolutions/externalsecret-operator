# External Secret Operator
![github actions](https://github.com/ContainerSolutions/externalsecret-operator/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ContainerSolutions/externalsecret-operator)](https://goreportcard.com/report/github.com/ContainerSolutions/externalsecret-operator) [![codecov](https://codecov.io/gh/ContainerSolutions/externalsecret-operator/branch/master/graph/badge.svg)](https://codecov.io/gh/ContainerSolutions/externalsecret-operator)

This operator reads information from a third party service
like [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) and automatically injects the values as [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

# Table of Contents

 * [Quick start](#quick-start)
 * [Manifests](#manifests)
 * [What does it do?](#what-does-it-do)
 * [Architecture](#architecture)
 * [Secrets Backends](#secrets-backends)
    * [1Password](#1password)
        * [Prerequisites](#prerequisites)
        * [Integration Test](#integration-test)
        * [Operator Deployment](#operator-deployment)
    * [GCP/Google Secrets Manager](#gcpgoogle-secrets-manager)
        * [Prerequisites](#prerequisites)
        * [Deploying](#deploying)
 * [Contributing](#contributing)
 
## Quick start

<!-- If you want to jump right into action you can deploy the External Secrets Operator using the provided [helm chart](./deployments/helm/externalsecret-operator/README.md) or [manifests](./deploy). The following examples are specific to the AWS Secret Manager backend. -->

<!-- ### Helm

Here's how you can deploy the External Secret Operator in the `default` namespace.

```shell
export AWS_ACCESS_KEY_ID="AKIAYOURSECRETKEYID"
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_SECRET_ACCESS_KEY="OoXie5Mai6Qu3fakemeezoo4ahfoo6IHahch0rai"
helm upgrade --install asm1 --wait \
    --set operatorName="asm-example" \
    --set secret.data.Type="asm" \
    --set secret.data.Parameters.accessKeyID="$AWS_ACCESS_KEY_ID" \
    --set secret.data.Parameters.region="$AWS_DEFAULT_REGION" \
    --set secret.data.Parameters.secretAccessKey="$AWS_SECRET_ACCESS_KEY" \
    ./deployments/helm/externalsecret-operator/.
```

It will watch for `ExternalSecrets` with `Backend: asm-example` resources in the `default` namespace and it will inject a corresponding `Secret` with the value retrieved from AWS Secret Manager.

Look for more deployment options in the [README.md](./deployments/helm/externalsecret-operator/README.md) of the helm chart. -->

### Manifests
- Uncomment and update backend config to be used in `config/backend-config/kustomization.yaml` with valid values:

```yaml
resources:
# - backend-config-gsm.yaml
- backend-config-asm.yaml
# - backend-config-dummy.yaml
# - backend-config-onepassword.yaml
```

```yaml
%cat config/backend-config/backend-config-asm.yaml
...
operator-config.json: |-
  {
    "Type": "asm",
    "Parameters": {
      "accessKeyID": "AWS_ACCESS_KEY_ID",
      "region": "AWS_DEFAULT_REGION",
      "secretAccessKey": "AWS_SECRET_ACCESS_KEY"
    }
  }
```
<!-- The `deploy` target in the Makefile will substiute variables and deploy the
manifests for you. The following command will deploy the operator in the
`default` namespace:

```shell
export AWS_ACCESS_KEY_ID="AKIAYOURSECRETKEYID"
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_SECRET_ACCESS_KEY="OoXie5Mai6Qu3fakemeezoo4ahfoo6IHahch0rai"
export OPERATOR_NAME=asm-example
export BACKEND=asm
make deploy
```
It will watch for `ExternalSecrets` with `Backend: asm-example` resources in the `default` namespace and it will inject a corresponding `Secret` with the value retrieved from AWS Secret Manager. -->

## What does it do?

Given a secret defined in AWS Secrets Manager:

```shell
% aws secretsmanager create-secret \
  --name=example-externalsecret-key \
  --secret-string='this string is a secret'
```

and an `ExternalSecret` resource definition like this :

```yaml
% cat config/samples/secrets_v1alpha1_externalsecret.yaml
apiVersion: secrets.externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecret
metadata:
  name: externalsecret-sample
  namespace: system
spec:
  key: example-externalsecret-key
  backend: 36af4962.externalsecret-operator.container-solutions.com
  version: latest
```

The operator fetches the secret from AWS Secrets Manager and injects it as a
secret:

```shell
% make install
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
this string is a secret
```

## Architecture

In [this article](https://docs.google.com/document/d/1hA6eM0TbRYcsDybiHU4kFYIqkEmDFo5GWNzJ2N398cI) you can find more information about the architecture and design choices. 

Here's a high-level diagram of how things are put together.

![architecture](./assets/architecture.png)

## Secrets Backends

We would like to support as many backend as possible and it should be rather easy to write new ones. Currently supported or planned backends are:

* [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
* [1Password](https://1password.com/security/)
* [Keybase](https://keybase.io/)
* [Git-secret](https://git-secret.io/)
* [GCP/Google Secret Manager](https://cloud.google.com/secret-manager)

<!-- A contributing guide is coming soon! -->

### 1Password

#### Prerequisites

* An existing 1Password team account.
* A 1Password account specifically for the operator. Tip: Setup an email with the `+` convention: `john.doe+operator@example.org`
* Store the _secret key_, _master password_, _email_ and _url_ of the _operator_ account in your existing 1Password account. This screenshot shows which fields should be used to store this information.
* Our naming convention for the item account is 'External Secret Operator' concatenated with name of the Kubernetes cluster for instance 'External Secret Operator minikube'. This item name is also used for development.
  
![1Password operator account](https://raw.githubusercontent.com/containersolutions/externalsecret-operator/master/assets/1password-operator-account.png)

#### Integration Test 

The integration `secrets/onepassword/backend_integration_test.go` test checks whether a secret stored in 1Password can be read via the operator.

Create a secret in 1Password as follow. Create a vault called `test vault one`. Now add a new `Login` item with name `testkey`. Set its `password` field to `testvalue`. See the screenshot below.

![1Password secret](https://raw.githubusercontent.com/containersolutions/externalsecret-operator/master/assets/1password-secret.png)

To run the integration test do the following.

1. Sign in to your _existing_ 1password

```
$ eval $(op signin)
```

2. Set the `ITEM_VAULT` and `ITEM_NAME` environment variables to select the right 1Password item that contains credentials fo your _operator_ 1Password account.

```
$ export ITEM_NAME=External Secret Operator mykubernetescluster
$ export ITEM_VAULT=myvault
```

Now load the 1Password credentials of your _operator_ account into the environment

```
$ . deployments/source-onepassword-secrets.sh
```

Run the tests including the integration test with

```
$ go test -v ./pkg/onepassword/
```

#### Operator Deployment

Follow the steps below to deploy the operator.

1. Sign in to your _existing_ 1password

```
$ eval $(op signin)
```

2. Load the 1Password credentials of your _operator_ account into the environment

```
$ source config/scripts/source-onepassword-secrets.sh
```

4.  Deploy the operator

```
$ make deploy-onepassword
```

### GCP/Google Secrets Manager
#### Prerequisites
- Enabled and configured secret manager API on your GCP project. [Secret Manager Docs](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)

#### Deploying

- Uncomment and update backend config to be used in `config/backend-config/kustomization.yaml`:

```yaml
resources:
- backend-config-gsm.yaml
# - backend-config-asm.yaml
# - backend-config-dummy.yaml
# - backend-config-onepassword.yaml
```

- Update the gsm backend config `config/backend-config/backend-config-gsm.yaml` with values from the service account key

```yaml
%cat config/backend-config/backend-config-gsm.yaml
...
operator-config.json: |-
    {
      "Type": "gsm",
      "Parameters": {
        "projectID": "",
        "type": "",
        "privateKeyID": "",
        "privateKey": "",
        "clientEmail": "",
        "clientID": "",
        "authURI": "",
        "tokenURI": "",
        "authProviderX509CertURL": "",
        "clientX509CertURL": ""
      }
    }

```

-  Update the resource definition `config/samples/secrets_v1alpha1_externalsecret.yaml`
```yaml
% cat config/samples/secrets_v1alpha1_externalsecret.yaml
apiVersion: secrets.externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecret
metadata:
  name: externalsecret-sample
  namespace: system
spec:
  key: your-secret-key
  backend: 36af4962.externalsecret-operator.container-solutions.com
  version: your-secret-version
```

- The operator fetches the secret from GCP Secret Manager and injects it as a
secret:

```shell
% make install
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.your-secret-key}' | base64 -d
```

## Contributing

Yay! We welcome and encourage contributions to this project! 

See our [contributing document](./CONTRIBUTING.md) and
[Issues](https://github.com/ContainerSolutions/externalsecret-operator/issues) for
planned improvements and additions.
