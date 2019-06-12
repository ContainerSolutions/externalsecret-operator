#!/bin/bash

OP_ITEM_JSON=$(op get item "External Secret Operator - 1Password")

OP_DOMAIN=$(echo $OP_ITEM_JSON | jq -r '.overview.url')
OP_EMAIL=$(echo $OP_ITEM_JSON | jq -r '.details.sections[0].fields[] | select(.t == "email").v')
OP_SECRET_KEY=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "username").value')
OP_MASTER_PASSWORD=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "password").value')

export OPERATOR_CONFIG="{ \"Name\": \"onepassword\", \"Type\": \"onepassword\", \"Parameters\": {\"domain\": \"${OP_DOMAIN}\", \"email\": \"${OP_EMAIL}\", \"secretKey\": \"${OP_SECRET_KEY}\", \"masterPassword\": \"${OP_MASTER_PASSWORD}\", \"vault\": \"Personal\" }}"