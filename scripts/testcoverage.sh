#!/bin/bash
set -e

MODE="count"

echo "mode: $MODE" > profile.cov

for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -not -path "*/vendor/*" -not -path "resources" -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -covermode=$MODE -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
		rm $dir/profile.tmp
    fi
fi
done

go tool cover -func profile.cov

gocov convert profile.cov | gocov-xml > coverage.xml
