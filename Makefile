DOCKER_IMAGE ?= containersol/externalsecret-operator

NAMESPACE ?= "default"
BACKEND ?= "asm"
OPERATOR_NAME ?= "asm-example"

GIT_HASH	:= $(shell git rev-parse --short HEAD)
GIT_BRANCH 	:= $(shell git rev-parse --abbrev-ref HEAD | sed 's/\//-/')
GIT_TAG 	:= $(shell git describe --tags --abbrev=0 --always)
DOCKER_TAG 	:= $(shell ./build/scripts/determine_docker_tag.sh $(GIT_HASH) $(GIT_BRANCH) $(GIT_TAG))

.PHONY: build
build: operator-sdk
	./operator-sdk build $(DOCKER_IMAGE)

.PHONY: push
.EXPORT_ALL_VARIABLES: push
push: build
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(GIT_HASH)
	docker push $(DOCKER_IMAGE):$(GIT_HASH)

	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: deploy
.EXPORT_ALL_VARIABLES: deploy
deploy:
	kubectl apply -n $(NAMESPACE) -f ./deploy/service_account.yaml
	kubectl apply -n $(NAMESPACE) -f ./deploy/role.yaml
	envsubst < ./deploy/role_binding.yaml | kubectl apply -n $(NAMESPACE) -f  -
	kubectl apply -n $(NAMESPACE) -f ./deploy/crds/externalsecret-operator_v1alpha1_externalsecret_crd.yaml
	envsubst < deploy/secret-${BACKEND}.yaml | kubectl apply -n $(NAMESPACE) -f -
	envsubst < deploy/deployment.yaml | kubectl apply -n $(NAMESPACE) -f -

.PHONY: apply-onepassword
OPERATOR_NAME=onepassword
BACKEND=onepassword
.EXPORT_ALL_VARIABLES: apply-onepassword
apply-onepassword:
	@echo "Deploying service account..."
	@kubectl apply -n $(NAMESPACE) -f ./deploy/service_account.yaml
	@echo "Deploying role..."
	@kubectl apply -n $(NAMESPACE) -f ./deploy/role.yaml
	@echo "Deploying rolebinding..."
	@envsubst < ./deploy/role_binding.yaml | kubectl apply -n $(NAMESPACE) -f  -
	@echo "Deploying external operator CRD..."
	@kubectl apply -n $(NAMESPACE) -f ./deploy/crds/externalsecret-operator_v1alpha1_externalsecret_crd.yaml
	@echo "Deploying 1password operator config secret..."
	@envsubst < deploy/secret-${BACKEND}.yaml | kubectl apply -n $(NAMESPACE) -f -
	@echo "Deploying operator deployment..."
	@envsubst < deploy/deployment.yaml | kubectl apply -n $(NAMESPACE) -f -

.PHONY: delete-onepassword
.EXPORT_ALL_VARIABLES: delete-onepassword
delete-onepassword:
	@echo "Deleting 1password operator config secret..."
	kubectl delete secret externalsecret-operator-config
	@echo "Deleting operator deployment..."
	kubectl delete deployment externalsecret-operator

.PHONY: deploy-onepassword
.EXPORT_ALL_VARIABLES: deploy-onepassword
deploy-onepassword: push apply-onepassword
	
.PHONY: test
test:
	go test -v -short ./... -count=1

.PHONY: coverage
# include only code we write in coverage report, not generated
COVERAGE := ./pkg/controller/externalsecret... ./secrets/...
coverage:
	go test -short -race -coverprofile=coverage.txt -covermode=atomic $(COVERAGE)
	curl -s https://codecov.io/bash | bash

.PHONY: test-helm
RELEASE := test$(shell echo $$$$)
OPERATOR_NAME=$(RELEASE)
BACKEND=dummy
.EXPORT_ALL_VARIABLES: test-helm
test-helm:
	helm upgrade --install --wait $(RELEASE) \
		--set test.create=true \
		./deploy/helm
	helm test  $(RELEASE)
	helm uninstall $(RELEASE)

PLATFORM := $(shell bash -c '[ "$$(uname -s)" = "Linux" ] && echo linux-gnu || echo apple-darwin')
OPERATOR_SDK_VERSION := v0.9.0
OPERATOR_SDK_URL := https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-$(PLATFORM)
operator-sdk:
	curl -LJ -o $@ $(OPERATOR_SDK_URL)
	chmod +x $@
