## Run tests and generate quality reports
test: build-test
	${DOCKER} run -ti --rm \
		-p ${SERVICE_PORT}:${SERVICE_PORT} \
		-v $$(pwd)/${REPORT_ARTIFACTS}:/app/${REPORT_ARTIFACTS} \
		--env BRANCH \
		--name ${APPNAME}-test \
		${DOCKER_IMAGE}:test ${TEST_CMD:-make test-int}
	[[ "${TEST_CMD}" =~ coverhtml ]] && ${DOCKER} cp ${APPNAME}-test:/app/cover.html ./cover.html && open ./cover.html || true

## Build test docker image
build-test:
	${DOCKER} build \
		-t ${DOCKER_IMAGE}:test \
		-f docker/dockerfile.test \
		.

.PHONY: test

## Run tests and output coverage reports
cover: test-cover-int

## Run tests and open report on default web browser
coverhtml: test-coverhtml-int

## Run code linter and output report as text
checkstyle: test-checkstyle-int

# Internal targets are run on the test docker container,
# they are not intended to be run directly

cover-int:
	@scripts/commands/test_cover.sh cli

coverhtml-int:
	@scripts/commands/test_cover.sh html

checkstyle-int:
	@scripts/commands/test_style.sh display

test-int:
	@echoHeader "Running Tests"
	@scripts/commands/test_style.sh
	@scripts/commands/test_cover.sh

test-%:
	make TEST_CMD="make $*" test
