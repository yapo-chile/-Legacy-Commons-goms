#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoTitle "Checking format with gofmt:"
ERRORS=0
for file in ${GO_FILES}; do
	echo "checking ${file}"
	errors=$(gofmt -d -e -s ${file} | grep -c -E "^\+[^\+]")
	if [ "${errors}" -gt 0 ]; then
		ERRORS=$((ERRORS + errors))
		echoError " ... found $$ errors issues"
	else
		echo " ... ok"
	fi
done
if [ $ERRORS -gt 0 ]; then
	echoError "found ${ERRORS} errors"
	#exit 1
fi

echoTitle "Checking code issues with go vet"
go tool vet -v ./pkg 2>&1
OUTPUT=$?
if [ "${OUTPUT}" -eq 0 ]; then
	echo "No errors found!"
fi

echoTitle "Checking code issues with golint"
${GOLINT} 2>&1

echoTitle "Checking that the files don't have huge functions"
gocyclo -over 19 ${GO_FILES}               # forbid code with huge functions

echoTitle "Done"
