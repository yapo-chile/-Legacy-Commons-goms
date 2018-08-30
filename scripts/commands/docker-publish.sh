#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### CODE ##############

#Publishing is only allowed from Travis
if [[ -n "$TRAVIS" ]]; then
    echoTitle "Publishing docker image to Artifactory"
    container_cache login --username "${ARTIFACTORY_USER}" --password "${ARTIFACTORY_PWD}" "${DOCKER_REGISTRY}"
    if [[ -n "$DOCKER_TAG" ]]; then
        container_cache push "${DOCKER_IMAGE}:${DOCKER_TAG}"
    else
        container_cache push "${DOCKER_IMAGE}"
    fi
else
    echoError "DOCKER PUBLISHING IS ONLY ALLOWED IN TRAVIS"
fi
