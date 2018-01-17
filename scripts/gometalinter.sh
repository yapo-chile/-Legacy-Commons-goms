#!/bin/bash
set -e

COMMAND='gometalinter ./... --config ".gometalinter.json"'
if [[ $@ == **display** ]]; then
    COMMAND="${COMMAND}"
else
    COMMAND="${COMMAND} --checkstyle | tee /dev/tty > checkstyle-report.xml"
fi

eval ${COMMAND}
status=${PIPESTATUS[0]}

# We need to catch error codes that are bigger then 2,
# they signal that gometalinter exited because of underlying error.
if [ ${status} -ge 2 ]; then
    echo "gometalinter exited with code ${status}, check gometalinter errors"
    exit ${status}
fi

exit 0
