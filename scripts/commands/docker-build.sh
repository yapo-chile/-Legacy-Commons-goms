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

if [[ "$BUILD_BRANCH" != "" ]]; then
    export GIT_BRANCH=${BUILD_BRANCH}
fi

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${GIT_BRANCH}"
echo "GIT TAG: ${GIT_TAG}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "GIT COMMIT SHORT: ${GIT_COMMIT_SHORT}"
echo "BUILD CREATOR: ${BUILD_CREATOR}"

export GIT_BRANCH_LOWERCASE=$(echo "${GIT_BRANCH}" | awk '{print tolower($0)}'| sed 's/\//_/;')
DOCKER_ARGS=" --no-cache"

if [[ "$GIT_TAG" != "" ]]; then
     DOCKER_ARGS="${DOCKER_ARGS} \
         -t ${DOCKER_IMAGE}:${GIT_TAG}"
elif [[ "${GIT_BRANCH_LOWERCASE}" == "master" ]]; then
     DOCKER_ARGS="${DOCKER_ARGS} \
         -t ${DOCKER_IMAGE}:${GIT_BRANCH_LOWERCASE} -t ${DOCKER_IMAGE}:latest"
else
     DOCKER_ARGS="${DOCKER_ARGS} \
         -t ${DOCKER_IMAGE}:${GIT_BRANCH_LOWERCASE}"
fi

DOCKER_ARGS=" ${DOCKER_ARGS} \
    --build-arg GIT_BRANCH="$GIT_BRANCH" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_CREATOR="$BUILD_CREATOR" \
    --build-arg VERSION="$VERSION" \
    --build-arg APPNAME="$APPNAME" \
    --build-arg BINARY="${DOCKER_BINARY}" \
    --build-arg EXPOSE_PORT="${DOCKER_PORT}" \
    --build-arg SYSLOG_ENABLED="${SYSLOG_ENABLED}" \
    --build-arg SYSLOG_IDENTITY="${SYSLOG_IDENTITY}" \
    --build-arg LOG_LEVEL="${LOG_LEVEL}" \
    --build-arg STDLOG_ENABLED="${STDLOG_ENABLED}" \
    -f docker/dockerfile \
    ."

echo "args: ${DOCKER_ARGS}"
set -x
docker build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
