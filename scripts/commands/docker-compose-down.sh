#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoTitle "Stopping Docker app"


if [[ $1 == "dev" ]]; then
    echoTitle "Stopping container for development"
    docker-compose -f docker/docker-compose-dev.yml -p ${APPNAME} down
elif [[ "$1" == "prod" ]]; then
    echoTitle "Stopping container for production"
    docker-compose -f docker/docker-compose.yml -p ${APPNAME} down
else
    echoError "Option not supported. Options: dev, prod"
    exit 1
fi
