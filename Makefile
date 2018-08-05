DOCKER_IMAGE ?= containersol/externalconfig-operator
DOCKER_TAG ?= $(shell grep -Po 'Version = "\K.*?(?=")' version/version.go)

# export these if you want to use AWS secrets manager
AWS_ACCESS_KEY_ID ?= AKIACONFIGUREME
AWS_SECRET_ACCESS_KEY ?= Secretsecretconfigureme 
AWS_REGION ?= eu-west-1

.PHONY: build
build: deps
	operator-sdk build $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: minikube
.EXPORT_ALL_VARIABLES: minikube
minikube: 
	make build
	kubectl apply -f ./deploy/rbac.yaml
	kubectl apply -f ./deploy/crd.yaml
	envsubst < deploy/operator-aws.yaml | kubectl apply -f -

.PHONY: test
test: deps
	go test -v ./...

.PHONY: deps
deps:
	dep ensure

.PHONY: install
install: deps
	go install -v ./...
	
