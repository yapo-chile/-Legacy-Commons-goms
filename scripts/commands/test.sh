#!/usr/bin/env bash

DIR="${BASH_SOURCE%/*}"

echoHeader "Running Tests"

"$DIR/test_style.sh"
"$DIR/test_cover.sh"
