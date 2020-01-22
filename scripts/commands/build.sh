#!/usr/bin/env bash

echoTitle "Building code"
set -e

go build -v -o ${APPNAME} ./${MAIN_FILE}

set +e
echoTitle "Done"
