#!/bin/bash

OP_ITEM_JSON=$(op get item "OnePassword - External Secret Operator")

export ONEPASSWORD_MASTER_PASSWORD=$(echo $OP_ITEM_JSON | jq -r .details.fields[0].value)
export ONEPASSWORD_EMAIL=$(echo $OP_ITEM_JSON | jq -r .details.sections[0].fields[0].v)
export ONEPASSWORD_SECRET_KEY=$(echo $OP_ITEM_JSON | jq -r .details.sections[0].fields[1].v)
export ONEPASSWORD_URL=$(echo $OP_ITEM_JSON | jq -r .overview.url)
