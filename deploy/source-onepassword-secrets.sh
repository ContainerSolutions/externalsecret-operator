#!/bin/bash

ITEM_NAME="${ITEM_NAME:-External Secret Operator minikube}"
ITEM_VAULT="${ITEM_VAULT:-Personal}"

OP_ITEM_JSON=$(op get item "${ITEM_NAME}" --vault=${ITEM_VAULT})

export OP_DOMAIN=$(echo $OP_ITEM_JSON | jq -r '.overview.url')
export OP_EMAIL=$(echo $OP_ITEM_JSON | jq -r '.details.sections[0].fields[] | select(.t == "email").v')
export OP_SECRET_KEY=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "username").value')
export OP_MASTER_PASSWORD=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "password").value')
export OP_VAULT="test vault one"
export BACKEND="onepassword"

export OPERATOR_CONFIG="{ \"Name\": \"onepassword\", \"Type\": \"onepassword\", \"Parameters\": {\"domain\": \"${OP_DOMAIN}\", \"email\": \"${OP_EMAIL}\", \"secretKey\": \"${OP_SECRET_KEY}\", \"masterPassword\": \"${OP_MASTER_PASSWORD}\", \"vault\": \"${OP_VAULT}\" }}"