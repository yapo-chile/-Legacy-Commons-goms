## Create docker image based on docker/dockerfile
docker-build: docker-boot
	@scripts/commands/docker-build.sh

## Push docker image to containers.mpi-internal.com
docker-publish: docker-build
	@scripts/commands/docker-publish.sh

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
