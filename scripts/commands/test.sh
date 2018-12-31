#!/usr/bin/env bash

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running Tests"

env GO111MODULE=off "$DIR/test_style.sh"
env GO111MODULE=on "$DIR/test_cover.sh"
