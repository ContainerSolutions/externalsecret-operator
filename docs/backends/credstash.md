## Credstash (AWS KMS)

#### Prerequisites

Create a KMS key in IAM, using an aws profile you have configured in the aws CLI. You can ommit --profile if you use the Default profile.

```
aws --region ap-southeast-2 --profile [yourawsprofile] kms create-key --query 'KeyMetadata.KeyId'
```

Assign the credstash alias to the key using the key id printed when you created the KMS key.

```
aws --region ap-southeast-2 --profile [yourawsprofile] kms create-alias --alias-name 'alias/credstash' --target-key-id "xxxx-xxxx-xxxx-xxx
```

Use a credstash client to create a secret (Using security context securityKey=securityValue here).

```
credstash put example-externalsecret-key  secretValue securityKey=securityValue
```


- Install CRDs 
```
  make install
```

#### Deployment

- Uncomment and update credentials to be used in `config/credentials/kustomization.yaml`:

```yaml
resources:
# - credentials-gsm.yaml
# - credentials-asm.yaml
# - credentials-dummy.yaml
# - credentials-gitlab.yaml
- credentials-credstash.yaml
```

- Update the credstash credentials `config/credentials/credentials-credstash.yaml` with correct AWS credentials.

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
    type: credstash
    auth: 
      secretRef: 
        name: externalsecret-operator-credentials-credstash
    parameters:
      region: eu-west-2
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

- The operator fetches the secret from AWS KMS like a credstash client
secret:

```shell
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
```