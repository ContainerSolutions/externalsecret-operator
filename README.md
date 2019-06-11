# External Secret Operator
[![CircleCI](https://circleci.com/gh/ContainerSolutions/externalsecret-operator.svg?style=svg)](https://circleci.com/gh/ContainerSolutions/externalsecret-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/ContainerSolutions/externalsecret-operator)](https://goreportcard.com/report/github.com/ContainerSolutions/externalsecret-operator) [![codecov](https://codecov.io/gh/ContainerSolutions/externalsecret-operator/branch/master/graph/badge.svg)](https://codecov.io/gh/ContainerSolutions/externalsecret-operator)

This operator reads information from a third party service
like [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) and automatically injects the values as [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

## Quick start

If you want to jump right into action you can deploy the External Secrets Operator using the provided [helm chart](./deploy/helm/README.md) or [manifests](./deploy/).

The following examples are specific to the AWS Secret Manager backend.

### Helm

The following command deploys the External Secret Operator in the `default` namespace.

```shell
export AWS_ACCESS_KEY_ID="AKIAYOURSECRETKEYID"
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_SECRET_ACCESS_KEY="OoXie5Mai6Qu3fakemeezoo4ahfoo6IHahch0rai"
helm upgrade --install asm1 --wait \
    --set secret.data.Name="asm-example" \
    --set secret.data.Type="asm" \
    --set secret.data.Parameters.accessKeyID="$AWS_ACCESS_KEY_ID" \
    --set secret.data.Parameters.region="$AWS_DEFAULT_REGION" \
    --set secret.data.Parameters.secretAccessKey="$AWS_SECRET_ACCESS_KEY" \
    ./deploy/helm/.
```

It will watch for `ExternalSecrets` with `Backend: asm-example` resources in the `default` namespace and it will inject a corresponding `Secret` with the value retrieved from AWS Secret Manager.

Look for more deployment options in the [README.md](./deploy/helm/README.md) of the helm chart.

### Manifests

For convenience this repository contains a set of deployment manifests in the `./deploy` directory. You can deploy them using the `deploy` make target:

```shell
export AWS_ACCESS_KEY_ID="AKIAYOURSECRETKEYID"
export AWS_DEFAULT_REGION="eu-west-1"
export AWS_SECRET_ACCESS_KEY="OoXie5Mai6Qu3fakemeezoo4ahfoo6IHahch0rai"
make deploy
```

The operator will be deployed in the `default` namespace and it will listen for `ExternalSecret` resources with `Backend: asm-example` in the whole cluster.

Check the manifests and the Makefile target for more details.

## What does it do?

Given a secret defined in AWS Secrets Manager:

```shell
% aws secretsmanager create-secret --name=example-externalsecret-key --secret-string='this string is a secret'
```

and an `ExternalSecret` resource definition like this one:

```yaml
% cat ./deploy/crds/examples/externalsecret-asm.yaml
apiVersion: externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecret
metadata:
  name: example-externalsecret
spec:
  Key: example-externalsecret-key
  Backend: asm
```

The operator fetches the secret from AWS Secrets Manager and injects it as a
secret:

```shell
% kubectl apply -f ./deploy/crds/examples/externalsecret-asm.yaml
% kubectl get secret example-externalsecret -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
this string is a secret
```

## Secrets Backends

We would like to support as many backend as possible and it should be rather easy to write new ones. Currently supported or planned backends are:

* AWS Secrets Manager
* One Password
* Keybase
* Git

A contributing guide is coming soon!

## What's next

This project is just at its beginning. See
[Issues](https://github.com/ContainerSolutions/externalsecret-operator/issues)
for planned improvements and additions.
