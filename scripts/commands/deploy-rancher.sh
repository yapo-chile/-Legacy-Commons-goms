#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

if [ -n "$1" ]; then
  ENV=$1
  echo "ENV targeted to $1"
else
  echo "Using default ENV ${RANCHER_ENV}"
fi

BRANCH=$(git branch | sed -n 's/^\* //p' | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')

RANCHER_CREDENTIALS=${RANCHER_API_KEY}:${RANCHER_SECRET_KEY}
RANCHER_STACK_NAME=${RANCHER_STACK_PREFIX}${ENV}

RANCHER_GET_STACK_URL="${RANCHER_API_URL}/v2-beta/projects/${RANCHER_INC_ID}/stacks?name=${RANCHER_STACK_NAME}"
RANCHER_STACK_ID=$(curl -s -u "${RANCHER_CREDENTIALS}" -X GET "${RANCHER_GET_STACK_URL}" | tee dump1.log | python -c "import sys,json; print json.load(sys.stdin)['data'][0]['id']")

RANCHER_GET_SERVICE_URL="${RANCHER_API_URL}/v2-beta/projects/${RANCHER_INC_ID}/services?stackId=${RANCHER_STACK_ID}&name=${APPNAME}"
RANCHER_SERVICE_ID=$(curl -s -u "${RANCHER_CREDENTIALS}" -X GET "${RANCHER_GET_SERVICE_URL}" | tee dump2.log | python -c "import sys,json; print json.load(sys.stdin)['data'][0]['id']")

echo RANCHER_STACK_ID $RANCHER_STACK_ID
echo RANCHER_SERVICE_ID $RANCHER_SERVICE_ID

RANCHER_UPGRADE_URL="${RANCHER_API_URL}/v2-beta/projects/${RANCHER_INC_ID}/services/${RANCHER_SERVICE_ID}/?action=upgrade"
RANCHER_FINISH_UPGRADE_URL="${RANCHER_API_URL}/v2-beta/projects/${RANCHER_INC_ID}/services/${RANCHER_SERVICE_ID}/?action=finishupgrade"
RANCHER_ROLLBACK_URL="${RANCHER_API_URL}/v2-beta/projects/${RANCHER_INC_ID}/services/${RANCHER_SERVICE_ID}/?action=rollback"

curl -u "${RANCHER_CREDENTIALS}" \
  -X POST \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"inServiceStrategy":{"launchConfig":{"tty":true, "vcpu":1, "imageUuid":"docker:'"${DOCKER_IMAGE}"':'"${BRANCH}"'", "startFirst":true}}, "toServiceStrategy":null}' \
  "${RANCHER_UPGRADE_URL}"
