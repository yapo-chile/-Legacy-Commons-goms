#!/usr/bin/env bash

#Build code again now for docker platform
echoHeader "Building production docker image"

set -x
${DOCKER} build \
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
	.
set +x

echoTitle "Build done"
