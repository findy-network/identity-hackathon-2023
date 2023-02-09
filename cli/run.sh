#!/bin/bash

set -e

ORIG_FCLI_JWT=$FCLI_JWT
ORIG_FCLI_CONFIG=$FCLI_CONFIG
ORIG_FCLI_USER=$FCLI_USER

CRED_DEF_ID=$(cat CRED_DEF_ID || echo "")

function createCredDef {
  if [ -z "$CRED_DEF_ID" ]; then
    echo "Create schema"
    sch_id=$(findy-agent-cli agent create-schema \
      --name="foobar" \
      --version=1.0 foo)

    # read schema - make sure it's found in ledger
    echo "Read schema $sch_id"
    schema=$(findy-agent-cli agent get-schema --schema-id $sch_id)

    # create cred def
    echo "Create cred def with schema id $sch_id"
    CRED_DEF_ID=$(findy-agent-cli agent create-cred-def \
      --id $sch_id --tag "$FCLI_USER")

    # read cred def - make sure it's found in ledger
    echo "Read cred def $CRED_DEF_ID"
    cred_def=$(findy-agent-cli agent get-cred-def --id $CRED_DEF_ID)

    echo $CRED_DEF_ID >"./CRED_DEF_ID"
  fi
}

function login {
  local jwt=$(findy-agent-cli authn login || echo "")
  if [ -z "$jwt" ]; then
    findy-agent-cli authn register
    jwt=$(findy-agent-cli authn login)
  fi
  export FCLI_JWT=$jwt
}

current_dir=$(dirname "$BASH_SOURCE")

if [ -z "$FCLI_URL" ]; then
  export FCLI_CONFIG="./config.yaml"
fi

if [ -z "$FCLI_USER" ]; then
  export FCLI_USER="cli-example"
fi

login
createCredDef

# replace cred def id in bot configs
sub_cmd='{sub("<CRED_DEF_ID>","'$CRED_DEF_ID'")}1'
awk "$sub_cmd" "issue-bot.template.yaml" >"issue-bot.yaml"
awk "$sub_cmd" "verify-bot.template.yaml" >"verify-bot.yaml"

conn_id=$(echo $(uuidgen) | tr '[:upper:]' '[:lower:]')
invitation=$(findy-agent-cli agent invitation --label $FCLI_USER -u --conn-id=$conn_id)

printf "\n\nHi there 👋 \n"
printf "\nPlease read the QR code with your wallet application to receive the credential.\n\n"

qrencode -m 2 -t utf8i <<<$invitation

printf "\n$invitation\n"

printf "\nIssue bot started 🤖\n"

findy-agent-cli bot start --conn-id $conn_id $current_dir/issue-bot.yaml

conn_id=$(echo $(uuidgen) | tr '[:upper:]' '[:lower:]')
invitation=$(findy-agent-cli agent invitation --label $FCLI_USER -u --conn-id=$conn_id)

printf "\n\nHi there 👋 \n"
printf "\nPlease read the QR code with your wallet application to prove your degree credential.\n\n"

qrencode -m 2 -t utf8i <<<$invitation

printf "\n$invitation\n"

printf "\nVerify bot started 🤖\n"

findy-agent-cli bot start --conn-id $conn_id $current_dir/verify-bot.yaml

export FCLI_JWT=$ORIG_FCLI_JWT
export FCLI_CONFIG=$ORIG_FCLI_CONFIG
export FCLI_USER=$ORIG_FCLI_USER
