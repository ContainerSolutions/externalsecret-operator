# External Config Operator

*This is a PoC*

I wanted to build an operator that reads information from a third party service
like [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/) or [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) and automatically inject the values as [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/).

That's pretty much what this operator does. This is in very early stage of
development and only AWS Secrets Manager is barely supported.

## Getting started

I used the [Operator SDK](https://github.com/operator-framework/operator-sdk)
to kickstart the project. It automatically uses
[dep](https://github.com/golang/dep) to handle dependencies (which I believe is
the right choice, anyway).

To build the project:
```
make build
```
This step will build a docker image and a simple deployment manifest.

The whole thing is working on
[minikube](https://github.com/kubernetes/minikube). You need to export your AWS
credentials so the operator can access AWS Secrets Manager and target the
minikube docker instance:

```
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
```
% aws secretsmanager get-secret-value --secret-id asecret --query SecretString
"secret"
```

and an `ExternalSecret` resource definition like this one:
```
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

```
% kubectl apply -f deploy/cr.yaml
% kubectl get secret asecret -o=go-template='{{ .data.asecret }}' | base64 -d
secret
```

## What's next
This could be just the beginning. If it seems like a good idea to continue
development there are many things to add, for example:
* more tests
* proper secrets/configuration backend configuration implementation
* more secrets/configuration backends
* helm chart to handle deployment
* a single ExternalSecret with a list of Secrets
* support ConfigMaps
