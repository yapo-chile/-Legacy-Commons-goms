## Build and start the service in development mode (detached)
run: mod build-dev "docker-compose-up -d"

## Build and start the service in development mode (attached)
start: mod build-dev docker-compose-up

.PHONY: run start

## Build develoment docker image
build-dev: docker-compose-build

## Setup directory for go module cache
mod:
	mkdir -p mod

## Run docker compose commands with the project configuration
docker-compose-%:
	docker-compose -f docker/docker-compose.yml \
		--project-name ${APPNAME} \
		--project-directory . \
		$*
