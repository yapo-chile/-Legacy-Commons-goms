#!/usr/bin/env bash

echoTitle "Building code"
set -e
echo "binary: ${PACT_BINARY} main: ${PACT_MAIN_FILE}"
go build -v -o ${PACT_BINARY} ./${PACT_MAIN_FILE}

set +e
echoTitle "Done"
