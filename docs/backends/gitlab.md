<a name="gitlab-cicd-variables"></a>

## Gitlab CI/CD Variables

<a name="gitlab-cicd-variables-pre"></a>

#### Prerequisites
- A Gitlab project with a CI/CD variable, who's key is `example_externalsecret_key`
- The project ID which you can find at the top of the main page of the project, right below the project name.
- A [Gitlab personal access token](https://gitlab.com/-/profile/personal_access_tokens) with `read_api` permissions

- Install CRDs
```
  make install
```

<a name="gitlab-cicd-variables-deployment"></a>

#### Deployment

- Uncomment and update credentials to be used in `config/credentials/kustomization.yaml`:

```yaml
resources:
# - credentials-gsm.yaml
# - credentials-asm.yaml
# - credentials-dummy.yaml
- credentials-gitlab.yaml
```

- Update the gitlab credentials `config/credentials/credentials-gitlab.yaml` with your personal access token

```yaml
%cat config/credentials/credentials-gitlab.yaml
...
credentials.json: |-
    {
      "token": "abcdef12345"
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
    type: gitlab
    auth:
      secretRef:
        name: externalsecret-operator-credentials-gitlab
    parameters:
      baseURL: https://gitlab.com
      projectID: 12345678
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
    - key: example_externalsecret_key
      version: latest
```

- The operator fetches the CI/CD variable from Gitlab and injects it as a secret:

```shell
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.example_externalsecret_key}' | base64 -d
```