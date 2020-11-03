```
apiVerson: store.externalsecret-operator.container-solutions.com/v1alpha1
kind: SecretStore 
metadata: {...}
spec:

  # Required
  # Unique name used to differenciate between different environments i.e production-aws, staging-aws, development-
  # NOTE: It should be unique for each store to avoid issues!
  controller: "dev"

  # Required
  store:
    # Sample store types
    # AWS Secrets Manager
    # store:
    #   type: asm
    #   auth: 
    #     secretRef: 
    #       name: externalsecret-operator-credentials-asm
    #   parameters:
    #     region: eu-west-2
    
    # GCP Secret Manager
    # store:
    #   type: gsm
    #   auth: 
    #     secretRef: 
    #       name: externalsecret-operator-credentials-gsm
    #   parameters:
    #     projectID: external-secrets-operator

    # Onepassword
    # store:
    #   type: onepassword
    #   auth: 
    #     secretRef: 
    #       name: externalsecret-operator-credentials-onepassword
    #   parameters:
    #     vault: Personal
    #     email: email@email-provider.com
    #     domain: domain.onepassword.com

status: {}
```