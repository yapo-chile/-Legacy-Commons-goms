#!/usr/bin/env bash
export UNAMESTR = $(uname)
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/ and pact/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))
# Retrieve process pid from ${PID} file
get_pid=$(cat ${PID} 2> /dev/null)
# Determine whether a process with ${EXEC} name is running with ${PID} file pid
is_running=$(ps --no-headers -o cmd -p $(call get_pid) 2> /dev/null | grep ${EXEC} && echo $?)

# GIT variables
export GIT_BRANCH=$(shell git branch | sed -n 's/^\* //p')
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_TAG=$(shell git tag -l --points-at HEAD | tr '\n' '_' | sed 's/_$$//;')
export GIT_COMMIT_SHORT=$(shell git rev-parse --short HEAD)
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)
export GIT_BRANCH_LOWERCASE=$(shell echo "${GIT_BRANCH}" | awk '{print tolower($0)}'| sed 's/\//_/;')

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `reports-publisher/config.json`
export REPORT_ARTIFACTS=reports

#APP variables
# This variables are for the use of your microservice. This variables must be updated each time you are creating a new microservice
export APPNAME=goms
export EXEC=./${APPNAME}
export VERSION=0.0.1
export YO=`whoami`
export LISTEN_PORT=$(call genport,2)
export SERVER_ROOT=${PWD}
export SERVERNAME=0.0.0.0
export BASE_URL="http://${SERVERNAME}:${LISTEN_PORT}"
export MAIN_FILE=cmd/${APPNAME}/main.go
export SYSLOG_ENABLED=false
export STDLOG_ENABLED=true
export LOG_LEVEL=0

#Pact test variables
export PACT_MAIN_FILE=cmd/${APPNAME}-pact/main.go
export PACT_BINARY=${APPNAME}-pact

#DOCKER variables
export DOCKER_REGISTRY=containers.schibsted.io
export DOCKER_IMAGE=${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_BINARY=${APPNAME}.docker
export DOCKER_PORT=8080
export DOCKER_PORT_DEV=$(call genport,1)
