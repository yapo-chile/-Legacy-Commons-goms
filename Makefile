include scripts/commands/vars.mk
include scripts/commands/colors.mk

# Image information
export APPNAME ?= $(shell basename `git rev-parse --show-toplevel`)
export BRANCH ?= $(shell git branch | sed -n 's/^\* //p')
export COMMIT ?= $(shell git rev-parse HEAD)
export COMMIT_DATE ?= $(shell TZ="America/Santiago" git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")
export CREATION_DATE ?= $(shell date -u '+%Y%m%d_%H%M%S')
export CREATOR ?= $(shell git log --format=format:%ae | head -n 1)

# Docker environment
export DOCKER_REGISTRY ?= containers.mpi-internal.com
export DOCKER_IMAGE ?= ${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_TAG ?= $(shell echo ${BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
export DOCKER ?= docker

# Golang environment
export GO111MODULE ?= on

## Run tests and generate quality reports
test:
	@scripts/commands/test.sh

## Run tests and output coverage reports
cover:
	@scripts/commands/test_cover.sh cli

## Run tests and open report on default web browser
coverhtml:
	@scripts/commands/test_cover.sh html

## Run gometalinter and output report as text
checkstyle:
	@scripts/commands/test_style.sh display

## Install golang system level dependencies
## Compile and build the executable file for pact tests
pact-build:
	scripts/commands/pact-build.sh

## Execute pact tests
pact-test: pact-build
	scripts/commands/pact-test.sh
	
## Compile the code
build:
	@scripts/commands/build.sh

## Execute the service
run:
	@env APP_PORT=${SERVICE_PORT} ./${APPNAME}

## Compile and start the service
start: build run

## New development workflow

run-dev: mod build-dev
	${DOCKER} run -ti \
		-v $$(pwd):/app \
		-v /var/empty:/app/mod \
		-v $$(pwd)/mod:/go/pkg/mod \
		-p ${SERVICE_PORT}:${SERVICE_PORT} \
		--env APPNAME \
		--env MAIN_FILE \
		${DOCKER_IMAGE}:${DOCKER_TAG}

build-dev:
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f docker/dockerfile.dev \
		--build-arg APPNAME \
		--build-arg MAIN_FILE \
		.

run-test: mod build-test
	${DOCKER} run -ti \
		-p ${SERVICE_PORT}:${SERVICE_PORT} \
		--env APPNAME \
		--env MAIN_FILE \
		${DOCKER_IMAGE}:test

build-test:
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:test \
		-f docker/dockerfile.test \
		--build-arg APPNAME \
		--build-arg MAIN_FILE \
		.

## Setup directory for module cache
mod:
	mkdir -p mod

## Compile and start the service using docker
docker-start: build docker-build docker-compose-up info

## Stop docker containers
docker-stop: docker-compose-down

## Setup a new service repository based on goms
clone:
	@scripts/commands/clone.sh

## Deploy project image on rancher
deploy-rancher:
	@scripts/commands/deploy-rancher.sh

deploy-k8s:
	@scripts/commands/deploy-k8s.sh

## Run gofmt to reindent source
fix-format:
	@scripts/commands/fix-format.sh

## Display basic service info
info:
	@echo "YO           : ${YO}"
	@echo "ServerRoot   : ${SERVER_ROOT}"
	@echo "API Base URL : ${BASE_URL}"
	@echo "Healthcheck  : curl ${BASE_URL}/healthcheck"

include docs.mk
include docker.mk
include help.mk
