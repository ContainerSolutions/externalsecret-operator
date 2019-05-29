#!/bin/bash

OP_ITEM_JSON=$(op get item "External Secret Operator - 1Password")

export ONEPASSWORD_SECRET_KEY=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "username").value')
export ONEPASSWORD_MASTER_PASSWORD=$(echo $OP_ITEM_JSON | jq -r '.details.fields[] | select(.designation == "password").value')
export ONEPASSWORD_DOMAIN=$(echo $OP_ITEM_JSON | jq -r '.overview.url')
export ONEPASSWORD_EMAIL=$(echo $OP_ITEM_JSON | jq -r '.details.sections[0].fields[] | select(.t == "email").v')