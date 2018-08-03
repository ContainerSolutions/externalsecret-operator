#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/ContainerSolutions/externalconfig-operator/pkg/generated \
github.com/ContainerSolutions/externalconfig-operator/pkg/apis \
externalconfig-operator:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
