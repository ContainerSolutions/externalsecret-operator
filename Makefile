DOCKER_IMAGE ?= containersol/externalsecret-operator

NAMESPACE ?= "default"
BACKEND ?= "asm"

.PHONY: build
build: operator-sdk
	./operator-sdk build $(DOCKER_IMAGE)

.PHONY: push
.EXPORT_ALL_VARIABLES: push
push:
	./build/scripts/push.sh

.PHONY: deploy
.EXPORT_ALL_VARIABLES: deploy
deploy:
	kubectl apply -n $(NAMESPACE) -f ./deploy/service_account.yaml
	kubectl apply -n $(NAMESPACE) -f ./deploy/role.yaml
	envsubst < ./deploy/role_binding.yaml | kubectl apply -n $(NAMESPACE) -f  -
	kubectl apply -n $(NAMESPACE) -f ./deploy/crds/externalsecret-operator_v1alpha1_externalsecret_crd.yaml
	envsubst < deploy/secret-${BACKEND}.yaml | kubectl apply -n $(NAMESPACE) -f -
	envsubst < deploy/deployment.yaml | kubectl apply -n $(NAMESPACE) -f -

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
test-helm:
	helm upgrade --install --wait $(RELEASE) \
		--set test.create=true \
		./deploy/helm
	helm test --cleanup $(RELEASE)
	helm delete --purge $(RELEASE)

PLATFORM := $(shell bash -c '[ "$$(uname -s)" = "Linux" ] && echo linux-gnu || echo apple-darwin')
OPERATOR_SDK_VERSION := v0.8.1
OPERATOR_SDK_URL := https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-$(PLATFORM)
operator-sdk:
	curl -LJ -o $@ $(OPERATOR_SDK_URL)
	chmod +x $@
