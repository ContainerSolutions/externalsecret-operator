# External Secret Operator
[![CircleCI](https://circleci.com/gh/ContainerSolutions/externalsecret-operator.svg?style=svg)](https://circleci.com/gh/ContainerSolutions/externalsecret-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/ContainerSolutions/externalsecret-operator)](https://goreportcard.com/report/github.com/ContainerSolutions/externalsecret-operator)

This operator reads information from a third party service
like [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) and automatically injects the values as [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

This project is in early stage of development and only AWS Secrets Manager is
supported.

## Getting started

This project was kickstarted using the [Operator
SDK](https://github.com/operator-framework/operator-sdk). It automatically
uses [dep](https://github.com/golang/dep) to handle dependencies.

To build the project:

```shell
make build
```

This step will build a docker image and a simple deployment manifest.

You need to export your AWS credentials so the operator can access AWS
Secrets Manager. Variables substitution is done by [`envsubst`](https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html) so you might need to install it.

This will build the project and deploy the operator and the required rbac roles
and custom resource definitions:


```shell
make deploy
```

## What does it do?

Given a secret defined in AWS Secrets Manager:

```shell
% aws secretsmanager create-secret --name=example-externalsecret-key --secret-string='this string is a secret'
```

and a `ExternalSecretBackend` resource as follows:

```yaml
% cat ./deploy/crds/examples/externalsecretbackend-asm.yaml
apiVersion: externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecretBackend
metadata:
  name: asm-example
spec:
  Type: asm
  Parameters:
    accessKeyID: AKIA...
    secretAccessKey: KSKSe4cret...
    region: eu-west-1
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
% kubectl apply -f ./deploy/crds/examples/externalsecretbackend-asm.yaml
% kubectl apply -f ./deploy/crds/examples/externalsecret-asm.yaml
% kubectl get secret example-externalsecret -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
this string is a secret
```

## What's next

This project is just at its beginning. See
[Issues](https://github.com/ContainerSolutions/externalsecret-operator/issues)
for planned improvements and additions.
