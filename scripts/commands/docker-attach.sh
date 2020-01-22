 #!/usr/bin/env bash

echoHeader "Attaching to docker container"
set -e
CONTAINER_ID=$(docker ps |grep ${DOCKER_IMAGE} | awk '{print $1}')

if [[ $CONTAINER_ID == "" ]]; then
    echoTitle "Docker Image not started. Please start with make docker-compose-up"
    exit 0
fi

docker attach --sig-proxy=false $CONTAINER_ID

set +e
