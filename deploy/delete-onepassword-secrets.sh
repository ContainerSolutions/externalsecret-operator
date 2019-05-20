#!/bin/bash

kubectl delete -n external-secret-operator secret onepassword-master-password
kubectl delete -n external-secret-operator secret onepassword-email
kubectl delete -n external-secret-operator secret onepassword-url
kubectl delete -n external-secret-operator secret onepassword-access-key
