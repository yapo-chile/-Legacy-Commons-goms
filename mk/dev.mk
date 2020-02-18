## Compile and start the service
run start: run-dev

.PHONY: run start

## Run the service in development mode
run-dev: mod build-dev
	${DOCKER} run -ti --rm \
		-v $$(pwd):/app \
		-v /var/empty:/app/mod \
		-v $$(pwd)/mod:/go/pkg/mod \
		-p ${SERVICE_PORT}:${SERVICE_PORT} \
		--env APPNAME \
		--name ${APPNAME} \
		${DOCKER_IMAGE}:${DOCKER_TAG}

## Build develoment docker image
build-dev:
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		-f docker/dockerfile.dev \
		--build-arg APPNAME \
		.

## Setup directory for go module cache
mod:
	mkdir -p mod
