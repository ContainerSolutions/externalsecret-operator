# External Secret Operator
[![CircleCI](https://circleci.com/gh/ContainerSolutions/externalsecret-operator.svg?style=svg)](https://circleci.com/gh/ContainerSolutions/externalsecret-operator) [![Go Report Card](https://goreportcard.com/badge/github.com/ContainerSolutions/externalsecret-operator)](https://goreportcard.com/report/github.com/ContainerSolutions/externalsecret-operator)

This operator reads information from a third party service
like [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) and automatically injects the values as [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

This project is in early stage of development and only AWS Secrets Manager is
supported.

## Getting started

This project was kickstarted using the [Operator
SDK](https://github.com/operator-framework/operator-sdk). It automatically uses
[dep](https://github.com/golang/dep) to handle dependencies.

To build the project:

```shell
make build
```

This step will build a docker image and a simple deployment manifest.

The whole thing is working on
[minikube](https://github.com/kubernetes/minikube). You need to export your AWS
credentials so the operator can access AWS Secrets Manager and target the
minikube docker instance:

```shell
eval $(minikube docker-env)
export AWS_ACCESS_KEY_ID=AKIACONFIGUREME
export AWS_SECRET_ACCESS_KEY=Secretsecretconfigureme 
export AWS_REGION=eu-west-1
make minikube
```
This will build the project and deploy the operator and the required rbac roles
and custom resource definitions.

## What does it do?
Given a secret defined in AWS Secrets Manager:

```shell
% aws secretsmanager get-secret-value --secret-id asecret --query SecretString
"secret"
```

and an `ExternalSecret` resource definition like this one:

```yaml
% cat deploy/cr.yaml 
apiVersion: "externalsecret-operator.container-solutions.com/v1alpha1"
kind: "ExternalSecret"
metadata:
  name: "asecret"
spec:
  Key: "asecret"
  Backend: "asm"
```

The operator fetches the secret from AWS Secrets Manager and injects it as a
secret:

```shell
% kubectl apply -f deploy/cr.yaml
% kubectl get secret asecret -o=go-template='{{ .data.asecret }}' | base64 -d
secret
```

## What's next

This project is just at its beginning. See
[Issues](https://github.com/ContainerSolutions/externalsecret-operator/issues)
for planned improvements and additions.