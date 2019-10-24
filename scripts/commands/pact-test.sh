#!/usr/bin/env bash
export PACT_TEST_ENABLED=true
export PACT_DIRECTORY=./pact
export SEARCH_MS_PORT=8089
export AD_API_PATH=http://localhost:${SEARCH_MS_PORT}/api/v1
export AD_GET_INFO=${AD_API_PATH}/search/list_id
export TOUCHSTONE_PORT=8088
export TOUCHSTONE_API_PATH=http://localhost:${TOUCHSTONE_PORT}
export TOUCHSTONE_GET_INFO=${TOUCHSTONE_API_PATH}/suggestions/marketplaces/yapocl/ads/
export SERVICE_HOST=:8090

file=pact-go_$(uname -s)_amd64.tar.gz

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

# Validate pact-go binaries
if [ ! -f "$PACT_DIRECTORY/bin/pact-go" ]; then
  echoTitle "Downloading binaries..."
  mkdir -p ${PACT_DIRECTORY}/bin
  wget --quiet -P ${PACT_DIRECTORY}/bin https://github.com/pact-foundation/pact-go/releases/download/v0.0.13/${file}
  tar zxf ${PACT_DIRECTORY}/bin/${file} -C ${PACT_DIRECTORY}/bin/
fi

echoTitle "Starting pact-go daemon in background"
nohup pact/bin/pact-go daemon > daemon.out 2> daemon.err &
PACT_PID=$!

: ' echoTitle "Starting search-ms mock in background"
nohup pact/bin/pact/bin/pact-stub-service pact/mocks/search-ms.json --port=${SEARCH_MS_PORT} &
SEARCH_MS_PID=$!

echoTitle "Starting touchstone-api mock in background"
nohup pact/bin/pact/bin/pact-stub-service pact/mocks/touchstone.json --port=${TOUCHSTONE_PORT} &
TOUCHSTONE_PID=$! '

echoTitle "Starting ${PACT_BINARY} in background"
nohup ./${PACT_BINARY} > ${PACT_BINARY}.out 2> ${PACT_BINARY}.err &
MS_PID=$!

#sleep 10
cd pact
go test -v -run TestProvider

echoTitle "Killing daemons"
kill ${PACT_PID} ${SEARCH_MS_PID} ${TOUCHSTONE_PID} ${MS_PID}

echoTitle "Done"