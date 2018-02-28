#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running Tests"

set -e

mkdir -p ${REPORT_ARTIFACTS}

COVER_FILE=${REPORT_ARTIFACTS}/cover.out
TMP_COVER_FILE=${REPORT_ARTIFACTS}/cover.out.tmp
COVERAGE_REPORT=${REPORT_ARTIFACTS}/coverage.xml
JUNIT_REPORT=${REPORT_ARTIFACTS}/junit-report.xml

echoTitle "Running Unit Tests"
function run_tests {
	# Get packages list except vendor and pact directories
	packages=$(go list ./... | grep -v vendor | grep -v pact )
	# Create cover output file
	echo "mode: count" > ${COVER_FILE}
	# Test all packages from the list
	for package in ${packages}; do
        echo "" > ${TMP_COVER_FILE}
        go test -v -covermode="count" -coverprofile=${TMP_COVER_FILE} ${package} || status=$?
		if [ -f ${TMP_COVER_FILE} ]; then
            cat ${TMP_COVER_FILE} | grep -v "mode: count" >> ${COVER_FILE}
		fi
	done
    grep -Ev "^$" ${COVER_FILE} > ${TMP_COVER_FILE}
    cat ${TMP_COVER_FILE} > ${COVER_FILE}
	return ${status:-0}
}

# Generate tests report
echoTitle "Generating tests report"
run_tests | tee /dev/tty | go-junit-report > ${JUNIT_REPORT}; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

# Print code coverage details
echoTitle "Printing code coverage details"
go tool cover -func ${COVER_FILE}

# Generate coverage report
echoTitle "Generating coverage report"
gocov convert ${COVER_FILE} | gocov-xml  > ${COVERAGE_REPORT}; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

echoTitle "Done"
exit ${status:-0}
