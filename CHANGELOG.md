# 0.1.0
- Updated CRD 
    - ExternalSecret
    - SecretStore
- Support backends additions
    - gsm
# Fixes


#91:

- Operator SDK updated to 1.0.1.
- OPERATOR_NAME not used in new controller-runtime, using - - - LeaderElectionID instead
deploy/ folder replaced with config/ which is handled by kustomize
- CRD manifests in deploy/crds/ are now in config/crd/bases
- CR manifests in deploy/crds/ are now in config/samples
- Controller manifest deploy/operator.yaml is now in config/manager/manager.yaml
- RBAC manifests in deploy are now in config/rbac/
- Go updated to 1.15
- Added tests to the externalsecret-controller
- Helpers and options to work with webhooks
use multigroup: true incase we need to support more complex secrets
- Add config/backend-config to handle backend-config secrets
Previous operator code in /legacy

- #47 - Add Dockerfile build stage and use https://github.com/GoogleContainerTools/distroless - generated by operator-sdk
- #3 - Introduce support for GCP secret manager(gsm) https://cloud.google.com/secret-manager/docs/configuring-secret-manager

- #86 - Migrate to github actions from circle ci


- #105
- #97
- #7 -Secret Binary Support
- #92
- #24
- #106
- #42
- #29