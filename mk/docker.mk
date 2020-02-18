## Create production build
build: docker-boot
	@echoHeader "Building production docker image"
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
	${DOCKER} tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}
	set +x

.PHONY: build

## Compile and start the service using docker
docker-start: build docker-build docker-compose-up info

## Stop docker containers
docker-stop: docker-compose-down

## Push docker image to containers.mpi-internal.com
docker-publish: build
	@echoTitle "Publishing docker image to Artifactory"
	${DOCKER} login --username "${ARTIFACTORY_USER}" --password "${ARTIFACTORY_PWD}" "${DOCKER_REGISTRY}"
	${DOCKER} push "${DOCKER_IMAGE}"


## Attach to this service's currently running docker container output stream
docker-attach:
	@scripts/commands/docker-attach.sh

## Start all required docker containers for this service
docker-compose-up: docker-boot
	@scripts/commands/docker-compose-up.sh

## Stop all running docker containers for this service
docker-compose-down:
	@scripts/commands/docker-compose-down.sh

## Start docker daemon
docker-boot:
	if ! [[ -n "$$TRAVIS" ]]; then \
		echoTitle "Starting Docker Engine"; \
		if [[ $$OSTYPE == "darwin"* ]]; then \
			echoTitle "Starting Mac OSX Docker Daemon"; \
			scripts/commands/docker-start-macosx.sh; \
		elif [[ "$$OSTYPE" == "linux-gnu" ]]; then \
			echoTitle "Starting Linux Docker Daemon"; \
			sudo start-docker-daemon; \
		else \
			echoError "Platform not supported"; \
		fi \
	fi
