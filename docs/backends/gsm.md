<a name="google-secret-manager"></a>

## GCP Secret Manager

<a name="google-secret-manager-pre"></a>

#### Prerequisites
- Enabled and configured secret manager API on your GCP project. [Secret Manager Docs](https://cloud.google.com/secret-manager/docs/configuring-secret-manager)

- Install CRDs 
```
  make install
```

<a name="google-secret-manager-deployment"></a>

#### Deployment

- Uncomment and update credentials to be used in `config/credentials/kustomization.yaml`:

```yaml
resources:
- credentials-gsm.yaml
# - credentials-asm.yaml
# - credentials-dummy.yaml
# - credentials-gitlab.yaml
```

- Update the gsm credentials `config/credentials/credentials-gsm.yaml` with service account key JSON

```yaml
%cat config/credentials/credentials-gsm.yaml
...
credentials.json: |-
    {
      "type": "service_account"
      ....
    }

```
-  Update the `SecretStore` resource definition `config/samples/store_v1alpha1_secretstore.yaml`
```yaml
% cat  `config/samples/store_v1alpha1_secretstore.yaml
apiVersion: store.externalsecret-operator.container-solutions.com/v1alpha1
kind: SecretStore
metadata:
  name: secretstore-sample
spec:
  controller: staging
  store:
    type: gsm
    auth: 
      secretRef: 
        name: externalsecret-operator-credentials-gsm
    parameters:
      projectID: external-secrets-operator
```

-  Update the `ExternalSecret` resource definition `config/samples/secrets_v1alpha1_externalsecret.yaml`
```yaml
% cat config/samples/secrets_v1alpha1_externalsecret.yaml
apiVersion: secrets.externalsecret-operator.container-solutions.com/v1alpha1
kind: ExternalSecret
metadata:
  name: externalsecret-sample
spec:
  storeRef: 
    name: externalsecret-operator-secretstore-sample
  data:
    - key: example-externalsecret-key
      version: latest
```

- The operator fetches the secret from GCP Secret Manager and injects it as a
secret:

```shell
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
```