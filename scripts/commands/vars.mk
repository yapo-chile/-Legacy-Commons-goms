#!/usr/bin/env bash

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `https://github.mpi-internal.com/spt-engprod/reports-publisher/blob/master/default_configuration.json`
export REPORT_ARTIFACTS=reports

# Pact test variables
export PACT_MAIN_FILE=cmd/${APPNAME}/main.go
export PACT_BINARY=${APPNAME}-pact
export PACT_DIRECTORY=pact
export PACT_TEST_ENABLED=false

#Pact broker
export PACT_BROKER_HOST=http://3.229.36.112
export PACT_BROKER_PORT=80
export PROVIDER_HOST=http://localhost
export PROVIDER_PORT=8080
export PACTS_PATH=./pacts

# Goms Client variables
export GOMS_HEALTH_PATH=${BASE_URL}/healthcheck

# User config
export PROFILE_HOST=http://10.15.1.78:7987
