#!/bin/bash

set -e

FCLI_PREV_CONFIG=$FCLI_CONFIG
FCLI_PREV_JWT=$FCLI_JWT

# make sure server is running
NOW=${SECONDS}
SERVER_WAIT_TIME=0
while ((${SERVER_WAIT_TIME} <= 60)); do
    printf "."
    SERVER_WAIT_TIME=$(($SECONDS - $NOW))
    if ((${SERVER_WAIT_TIME} >= 60)); then
        printf "\nServer start failed\n"
        exit 1
    fi
    RES_CODE=$(curl -s --write-out '%{http_code}' --output /dev/null http://localhost:3001)
    if ((${RES_CODE} == 200)); then
        SERVER_WAIT_TIME=61
    else
        sleep 1
    fi
done

# use yaml file for cli configuration
CURRENT_DIR=$(dirname "$BASH_SOURCE")
export FCLI_CONFIG=$CURRENT_DIR/config.yaml

findy-agent-cli authn register || echo "Ignoring error: probably already registered"
LOGIN_OUTPUT=$(findy-agent-cli authn login)
WORDS=($LOGIN_OUTPUT)
INDEX=$((${#WORDS[@]} - 1))
# token is the last word of returned string
export FCLI_JWT=${WORDS[$INDEX]}

# set automatic accept mode
findy-agent-cli agent mode-cmd -a

# start listening events
findy-agent-cli agent listen >test.log &

echo "*******************************************"
echo "Test 1: issue credential"

# fetch issue page and parse invitation
ISSUE_HTML=$(curl -s http://localhost:3001/issue)
INVITATION=$(echo $ISSUE_HTML | awk -v FS="(cols=\"60\">|</textarea>)" '{print $2}')

# make connection to service agent
CONNECTION_ID=$(findy-agent-cli agent connect --invitation $INVITATION)

# make sure issuing succeeds by checking status updates
NOW=${SECONDS}
ISSUE_WAIT_TIME=0
while ((${ISSUE_WAIT_TIME} <= 60)); do
    printf "."
    ISSUE_WAIT_TIME=$(($SECONDS - $NOW))
    if ((${ISSUE_WAIT_TIME} >= 60)); then
        printf "\nIssuing failed\n"
        exit 1
    fi
    if grep -q '| ISSUE_CREDENTIAL | STATUS_UPDATE |' test.log; then
        ISSUE_WAIT_TIME=61
    else
        sleep 1
    fi
done

echo "Test 1: issue credential SUCCESS"

echo "*******************************************"
echo "Test 2: verify proof"

# fetch verify page and parse invitation
VERIFY_HTML=$(curl -s http://localhost:3001/verify)
INVITATION=$(echo $VERIFY_HTML | awk -v FS="(cols=\"60\">|</textarea>)" '{print $2}')

# make connection to service agent
CONNECTION_ID=$(findy-agent-cli agent connect --invitation $INVITATION)

# make sure verifying succeeds by checking status updates
NOW=${SECONDS}
VERIFY_WAIT_TIME=0
while ((${VERIFY_WAIT_TIME} <= 60)); do
    printf "."
    VERIFY_WAIT_TIME=$(($SECONDS - $NOW))
    if ((${VERIFY_WAIT_TIME} >= 60)); then
        printf "\nVerifying failed\n"
        exit 1
    fi
    if grep -q '| PRESENT_PROOF | STATUS_UPDATE |' test.log; then
        VERIFY_WAIT_TIME=61
    else
        sleep 1
    fi
done

echo "Test 2: verify proof SUCCESS"
echo "*******************************************"

# reset env variables
export FCLI_CONFIG=$FCLI_PREV_CONFIG
export FCLI_JWT=$FCLI_PREV_JWT

# stop listener
kill -9 $(ps aux | pgrep -f 'findy-agent-cli agent listen') &>/dev/null
