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
    github.com/Masterminds/glide
    gopkg.in/alecthomas/gometalinter.v2
)

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
    env GO111MODULE=off go get -u -v ${tool}
done

echoTitle "Installing linters"
# Install all available linters
env GO111MODULE=off gometalinter.v2 --install

set +e
