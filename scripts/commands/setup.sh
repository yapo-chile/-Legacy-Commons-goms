#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"

set -e
# List of tools used for testing, validation, and report generation
tools=(
	github.com/golang/lint/golint
    github.com/axw/gocov/gocov
	github.com/AlekSi/gocov-xml
    gopkg.in/alecthomas/gometalinter.v1
    github.com/jstemmer/go-junit-report
    github.com/fzipp/gocyclo                             # Function cyclomatic complexity analyzer
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
	which $(basename ${tool}) > /dev/null || go get -u -v ${tool}
done

echoTitle "Installing linters"
# Install all available linters
gometalinter.v1 --install

echoTitle "Installing Glide dependencies"
glide install

set +e
