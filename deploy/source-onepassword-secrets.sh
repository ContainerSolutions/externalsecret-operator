#!/bin/bash

OP_ITEM_JSON=$(op get item "External Secret Operator - 1Password")

DOMAIN=$(echo $OP_ITEM_JSON | jq -r '.overview.url')
EMAIL=$(echo $OP_ITEM_JSON | jq -r '.details.sections[0].fields[] | select(.t == "email").v')
SECRET_KEY=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "username").value')
MASTER_PASSWORD=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "password").value')

export OPERATOR_CONFIG="{ \"Name\": \"onepassword\", \"Type\": \"onepassword\", \"Parameters\": {\"domain\": \"${DOMAIN}\", \"email\": \"${EMAIL}\", \"secretKey\": \"${SECRET_KEY}\", \"masterPassword\": \"${MASTER_PASSWORD}\" }}"