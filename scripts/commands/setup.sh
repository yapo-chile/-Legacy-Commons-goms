#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"

set -e

echoTitle "Removing outdated vendor"
rm -rf vendor glide.*

echoTitle "Initializating go modules"
GO111MODULE=on go mod init $APPMODULE

echoTitle "Installing project dependencies"
GO111MODULE=on go mod tidy

set +e
