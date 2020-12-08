## Azure Key Vault secrets

### Prerequisites

You need to have a Key Vault instance with a secret and an application registered in your Azure Active directory with Read access to the Vault's secrets.

The following script creates everything needed for sample purposes. It assumes you have Azure ClI installed and it is already authenticated.  
For more information about Azure CLI refer to the Azure CLI's [documentation page](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli).

```bash
# The tennatId will be needed to get the secrets
TENANT_ID=$(az account show --query tenantId | tr -d \")

# Name and location of the Resource Group
RESOURCE_GROUP="MyKeyVaultResourceGroup"
LOCATION="westus"

# Create the Resource Group
az group create --location $LOCATION --name $RESOURCE_GROUP

VAULT_NAME="eso-akv-test"

# Create the Key Vault
az keyvault create --name $VAULT_NAME --resource-group $RESOURCE_GROUP

SECRET_NAME="example-externalsecret-key"
SECRET_VAlUE="This is our secret now!"

# Add a secret to the vault
az keyvault secret set --name $SECRET_NAME --vault-name $VAULT_NAME --value "$SECRET_VAlUE"

# Now you need to create an app to access the Key Vault
APP_NAME="ExtSectret Query App"
APP_ID=$(az ad app create --display-name "$APP_NAME" --query appId | tr -d \")

# A Service Principal must also be created
SERVICE_PRINCIPAL=$(az ad sp create --id $APP_ID --query objectId | tr -d \")

# Add permission to your App to query the Key Vault
# The --api-permission refers to the Azure Key Vault user_impersonation permission (do not modify)
# The --api refers to the Azure Key Vault API (do not modify)
az ad app permission add --id $APP_ID --api-permissions f53da476-18e3-4152-8e01-aec403e6edc0=Scope --api cfa8b339-82a2-471a-a3c9-0fc0be7a4093

APP_PASSWORD="ThisisMyStrongPassword"
# A password must be created for the app
az ad app credential reset --id $APP_ID --password "$APP_PASSWORD"

# Finnaly, the Key Vault must have an Access Policy for the created app
az keyvault set-policy --name $VAULT_NAME --object-id $SERVICE_PRINCIPAL --secret-permissions get
```

For a detailed view on how to create the above mentioned resources in the Azure Portal, please go to: [How To Access Azure Key Vault Secrets Through Rest API Using Postman](https://www.c-sharpcorner.com/article/how-to-access-azure-key-vault-secrets-through-rest-api-using-postman/)

- Now you're ready to tnstall CRDs
```
  make install
```

### Deployment

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
    "tenantId": "<Active Directory's Tenant ID>",
    "clientId": "<Application (client) ID>",
    "clientSecret": "<Application's secret value>",
    "keyvault": "<Key Vault name>"
}
```

You can run the following script that will generate the above mentioned json object  
```bash
echo -e "{ \n \
  \"tenantId\": \"$TENANT_ID\", \n \
  \"clientId\": \"$APP_ID\", \n \
  \"clientSecret\": \"$APP_PASSWORD\", \n \
  \"keyvault\": \"$VAULT_NAME\" \n \
}"
```
> Beware of the indentation if you paste the output from above into your file.


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

### Clean up

- Delete the resource group (it willa lso delete the Kay Vault created)

```bash
az group delete --name $RESOURCE_GROUP 
```

- Delete the Active Directory application

```bash
az ad app delete --id $APP_ID
```