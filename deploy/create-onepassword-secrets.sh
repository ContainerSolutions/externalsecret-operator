#!/bin/bash

kubectl create -n external-secret-operator secret generic onepassword --from-literal=onepassword-master-password="${ONEPASSWORD_MASTER_PASSWORD}" \
                                                                      --from-literal=onepassword-email="${ONEPASSWORD_EMAIL}" \
                                                                      --from-literal=onepassword-domain="${ONEPASSWORD_DOMAIN}" \
                                                                      --from-literal=onepassword-secret-key="${ONEPASSWORD_SECRET_KEY}"
