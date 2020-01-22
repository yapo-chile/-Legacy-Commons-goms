#!/usr/bin/env bash

#Publishing is only allowed from Travis
if [[ -n "$TRAVIS" ]]; then
    echoTitle "Publishing docker image to Artifactory"
    ${DOCKER} login --username "${ARTIFACTORY_USER}" --password "${ARTIFACTORY_PWD}" "${DOCKER_REGISTRY}"
    ${DOCKER} push "${DOCKER_IMAGE}"
else
    echoError "DOCKER PUBLISHING IS ONLY ALLOWED IN TRAVIS"
fi
