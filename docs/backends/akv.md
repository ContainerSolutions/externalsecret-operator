# Azure Key Vault

This are the needed steps to use Azure Key Vault as a backend.
The External Secret Operator (ESO) will impersonate an application you must register in your Azure Active Directory.

Make sure you add the Azure Key Vault user_impersonation permission to the application. 

Also, under the application's Certificates & secrets, you need to create a Client Secret.

After your application is ready you need to create Access Policy in your Key Vault to allow the application access your secret. Just the 'Get' permission is needed under Secret Management Operations.

With that, you need to modify the credentials-akv.yaml file to add your credentials:

```json
{
    "tennant_id": "<Active Directory's Tennant ID>",
    "client_id": "<Application (client) ID>",
    "client_secret": "<Application's secret value>",
    "keyvault": "<Key Vault name>"
}
```

Ref: https://www.c-sharpcorner.com/article/how-to-access-azure-key-vault-secrets-through-rest-api-using-postman/