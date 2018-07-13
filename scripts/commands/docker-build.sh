#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### DYNAMIC VARS ###############

#In case we are in travis, docker tag will be "branch_name-20180101-1200". In case of master branch, branch-name is blank.
#In case of local build (not in travis) tag will be "local".
if [[ -n "$TRAVIS" ]]; then
    if [ "${GIT_BRANCH}" != "master" ]; then
        DOCKER_TAG=$(echo ${GIT_BRANCH}- | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')$(date -u '+%Y%m%d_%H%M%S')
    else
        DOCKER_TAG=$(date -u '+%Y%m%d_%H%M%S')
    fi
else
    DOCKER_TAG=local
fi


########### CODE ##############

#Build code again now for docker platform
echoHeader "Building code for docker platform"
set -e

rm -f ${DOCKER_BINARY}
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${DOCKER_BINARY} ./${MAIN_FILE}

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

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${GIT_BRANCH}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "BUILD CREATOR: ${BUILD_CREATOR}"
echo "IMAGE NAME: ${DOCKER_IMAGE}:${DOCKER_TAG}"

DOCKER_ARGS=" \
    -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
    --build-arg GIT_BRANCH="$GIT_BRANCH" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg BUILD_CREATOR="$BUILD_CREATOR" \
    --build-arg APPNAME="$APPNAME" \
    --build-arg BINARY="${DOCKER_BINARY}" \
    -f docker/dockerfile \
    ."

echo "args: ${DOCKER_ARGS}"
set -x
docker build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
