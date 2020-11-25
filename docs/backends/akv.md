## Azure Key Vault secrets

#### Prerequisites

- Create an Azure Key Vault resource with a secret who's name is `example-externalsecret-key`
- In your Azure Active Directory, create an app under App Registration and assign the 
- In your recently created app go to API Permissions and assign the `Azure Key Vault/user_impersonation` permission.
- Under Certificates & secretes create a client secret.
- Go back to your Key Vault resource and under Access policies, create an access policy for the app. Make sure you at least select 'Get' under 'Secret Management Operations'.
- To authenticate you will need: (to learn more about the authentication process take a look at [Authenticate to Azure Key Vault](https://docs.microsoft.com/en-us/azure/key-vault/general/authentication))
    - The Application (client) ID.
    - The Client Secret (also from the application)
    - The Tennant ID.
    - The Azure Key Vault service name.
- Install CRDs
```
  make install
```

For a detailed view on how to create the above mentioned resources, please go to: [How To Access Azure Key Vault Secrets Through Rest API Using Postman](Ref: https://www.c-sharpcorner.com/article/how-to-access-azure-key-vault-secrets-through-rest-api-using-postman/)


#### Deployment

- Uncomment and update credentials to be used in `config/credentials/kustomization.yaml`:

```yaml
resources:
# - credentials-gsm.yaml
# - credentials-asm.yaml
# - credentials-dummy.yaml
# - credentials-gitlab.yaml
- credentials-akv.yaml

```

- Update the Azure Key Vault backend credentials `config/credentials/credentials-akv.yaml` with your personal access token

```json
{
    "tennant_id": "<Active Directory's Tennant ID>",
    "client_id": "<Application (client) ID>",
    "client_secret": "<Application's secret value>",
    "keyvault": "<Key Vault name>"
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
    type: akv
    auth:
      secretRef:
        name: externalsecret-operator-credentials-akv
    parameters: {}
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
      version: ""
```

- The operator fetches the Key Vault secret from Azure and injects it as a Kubernetes secret:

```shell
% make deploy
% kubectl get secret externalsecret-operator-externalsecret-sample -n externalsecret-operator-system \
  -o jsonpath='{.data.example-externalsecret-key}' | base64 -d
```
