#!/bin/bash

kubectl create -n external-secret-operator secret generic onepassword-master-password --from-literal=onepassword_master_password="${ONEPASSWORD_MASTER_PASSWORD}"
kubectl create -n external-secret-operator secret generic onepassword-email --from-literal=onepassword_email="${ONEPASSWORD_EMAIL}"
kubectl create -n external-secret-operator secret generic onepassword-url --from-literal=onepassword_url="${ONEPASSWORD_URL}"
kubectl create -n external-secret-operator secret generic onepassword-access-key --from-literal=onepassword_master_key="${ONEPASSWORD_MASTER_KEY}"
