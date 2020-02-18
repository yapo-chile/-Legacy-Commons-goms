#!/usr/bin/env bash
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `reports-publisher/config.json`
export REPORT_ARTIFACTS=reports

# Pact test variables
export PACT_MAIN_FILE=cmd/${APPNAME}/main.go
export PACT_BINARY=${APPNAME}-pact
export PACT_DIRECTORY=pact
export PACT_TEST_ENABLED=false

# DOCKER variables
export DOCKER_PORT=$(call genport,1)

# Documentation variables
export DOCS_DIR=docs
export DOCS_HOST=localhost:$(call genport,3)
export DOCS_PATH=github.mpi-internal.com/Yapo/${APPNAME}
export DOCS_COMMIT_MESSAGE=Generate updated documentation

# Goms Client variables
export GOMS_HEALTH_PATH=${BASE_URL}/healthcheck

# User config
export PROFILE_HOST=http://10.15.1.78:7987

#Pact broker
export PACT_BROKER_HOST=http://3.229.36.112
export PACT_BROKER_PORT=80
export PROVIDER_HOST=http://localhost
export PROVIDER_PORT=8080
export PACTS_PATH=./pacts
