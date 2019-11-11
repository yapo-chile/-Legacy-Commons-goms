#!/usr/bin/env bash
export PACT_TEST_ENABLED=true
export PACT_DIRECTORY=./pact
export PROFILE_MS_PORT=7987
export SERVICE_HOST=:8090
export PROFILE_API_PATH=http://localhost:${SEARCH_MS_PORT}/api/v1

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

echoTitle "Starting profile-ms mock in background"
nohup pact/bin/pact/bin/pact-stub-service pact/mocks/profile-ms.json --port=${PROFILE_MS_PORT} &
PROFILE_PID=$!

echoTitle "Starting ${PACT_BINARY} in background"
nohup ./${PACT_BINARY} > ${PACT_BINARY}.out 2> ${PACT_BINARY}.err &
MS_PID=$!

#sleep 10
cd pact
go test -v -run TestProvider

echoTitle "Killing daemons"
kill ${PACT_PID} ${PROFILE_PID} ${MS_PID}

echoTitle "Done"