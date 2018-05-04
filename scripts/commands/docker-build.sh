#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

#Build code again now for docker platform
echoHeader "Building code for docker platform"
set -e

rm -f ${DOCKER_BINARY}
GOOS=linux GOARCH=386 go build -v -o ${DOCKER_BINARY} ./${MAIN_FILE}

set +e

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

if [ -n "${BUILD_BRANCH}" ]; then
    export GIT_BRANCH=${BUILD_BRANCH}
fi

BUILD_NAME=$(if [ -n "${GIT_TAG}" ]; then echo -n "${GIT_TAG}"; else echo -n "${GIT_BRANCH}"; fi)
export BUILD_TAG=$(echo -n "${BUILD_NAME}" | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${GIT_BRANCH}"
echo "GIT TAG: ${GIT_TAG}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "GIT COMMIT SHORT: ${GIT_COMMIT_SHORT}"
echo "BUILD CREATOR: ${BUILD_CREATOR}"
echo "BUILD NAME: ${DOCKER_IMAGE}:${BUILD_TAG}"

DOCKER_ARGS=""

if [[ "${GIT_BRANCH}" == "master" ]]; then
     DOCKER_ARGS="${DOCKER_ARGS} \
         -t ${DOCKER_IMAGE}:${BUILD_TAG} -t ${DOCKER_IMAGE}:latest"
else
     DOCKER_ARGS="${DOCKER_ARGS} \
         -t ${DOCKER_IMAGE}:${BUILD_TAG}"
fi

DOCKER_ARGS=" ${DOCKER_ARGS} \
    --build-arg GIT_BRANCH="$GIT_BRANCH" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_CREATOR="$BUILD_CREATOR" \
    --build-arg VERSION="$VERSION" \
    --build-arg APPNAME="$APPNAME" \
    --build-arg BINARY="${DOCKER_BINARY}" \
    -f docker/dockerfile \
    ."

echo "args: ${DOCKER_ARGS}"
set -x
docker build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
