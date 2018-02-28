#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoTitle "Starting Docker Engine"
if [[ $OSTYPE == "darwin"* ]]; then
    echoTitle "Starting Mac OSX Docker Daemon"
    $DIR/docker-start-macosx.sh
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
    echoTitle "Starting Linux Docker Daemon"
    sudo start-docker-daemon
else
    echoError "Platform not supported"
fi

echoTitle "Creating needed networks"
for network in ${DOCKER_COMPOSE_NETWORKS}; do
    networkId=`docker network ls -q -f name=${network}`
    if [ -z "$networkId" ];
    then
        echo "Creating network ${network}"
        docker network create ${network}
    fi
done

echoTitle "Starting docker compose"
if [[ $1 == "dev" ]]; then
    echoTitle "Starting container for development"
    docker-compose -f docker/docker-compose-dev.yml -p ${APPNAME} up -d
elif [[ "$1" == "prod" ]]; then
    echoTitle "Starting container for production"
    docker-compose -f docker/docker-compose.yml -p ${APPNAME} up -d
else
    echoError "Option not supported. Options: dev, prod"
    exit 1
fi

echoTitle "Done"
