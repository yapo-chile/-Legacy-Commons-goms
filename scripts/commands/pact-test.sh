#!/usr/bin/env bash

# Pact tests
export MS_MAIN_FILE=cmd/${APPNAME}/main.go
export MS_BINARY=${APPNAME}-pact

# Pact binaries and logs directories in container
export PACT_BINARY=../pact/bin
export PACT_LOGS=./reports

# Profile service
export PROFILE_MS_PORT=5555
export PROFILE_HOST=http://localhost:${PROFILE_MS_PORT}

echoTitle "Building microservice pact test binary"
go build -v -o ${MS_BINARY} ./${MS_MAIN_FILE}

echoTitle "Starting profile mock in background"
nohup ${PACT_BINARY}/pact-stub-service pact/mocks/profile.json --port=${PROFILE_MS_PORT} > ${PACT_LOGS}/profile.out 2>&1  &
PROFILE_PID=$!

echo ${PROFILE_PID}

echoTitle "Starting ${MS_BINARY} in background"
nohup  ./${MS_BINARY} > ${PACT_LOGS}/${MS_BINARY}.out 2> ${PACT_LOGS}/${MS_BINARY}.err &
MS_PID=$!

echo ${MS_PID}

sleep 10
cd pact
go test -v -run TestProvider
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]
then
  echoTitle "Error executing pact-tests"
elif [[ -n "$TRAVIS" ]]
then
  echoTitle "Publishing pact files into pact-broker."
  go test -v -run TestSendBroker
fi

echoTitle "Killing daemons"
kill -9 ${PROFILE_PID} ${MS_PID}

echoTitle "Done"
exit ${EXIT_CODE}
