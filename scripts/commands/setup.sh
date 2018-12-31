#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"

set -e
# List of tools used for testing, validation, and report generation
tools=(
    github.com/jstemmer/go-junit-report
    github.com/axw/gocov/gocov
    github.com/AlekSi/gocov-xml
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
    env GO111MODULE=off go get -u -v ${tool}
done

echoTitle "Installing linters"
# Install all available linters
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.12.5

set +e
