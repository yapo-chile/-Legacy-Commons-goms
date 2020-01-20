#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### DYNAMIC VARS ###############

#In case we are in travis, we will use cached docker environment.
if [[ -n "$TRAVIS" ]]; then
    DOCKER_COMMAND=container_cache
else
    DOCKER_COMMAND=docker
fi


########### CODE ##############

#Build code again now for docker platform
echoHeader "Building code for docker platform"

set +e

#In case it is not Travis we make sure docker is running
if ! [[ -n "$TRAVIS" ]]; then
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
fi

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${BRANCH}"
echo "GIT COMMIT: ${COMMIT}"
echo "GIT COMMIT DATE: ${COMMIT_DATE}"
echo "BUILD CREATOR: ${CREATOR}"
echo "BUILD DATE: ${CREATION_DATE}"
echo "IMAGE NAME: ${DOCKER_IMAGE}:${DOCKER_TAG}"

DOCKER_ARGS=" \
    -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
	-f docker/dockerfile \
	--build-arg APPNAME=${APPNAME} \
	--build-arg GIT_COMMIT=${COMMIT} \
	--label appname=${APPNAME} \
	--label branch=${BRANCH} \
	--label build-date=${CREATION_DATE} \
	--label commit=${COMMIT} \
	--label commit-author=${CREATOR} \
	--label commit-date=${COMMIT_DATE} \
	."

echo "args: ${DOCKER_ARGS}"
set -x
${DOCKER_COMMAND} build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
