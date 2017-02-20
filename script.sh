#!/bin/bash
ID=$(id -u)
BASEPORT=$(((ID - ($ID / 100 ) * 100) * 200 + 30000))
MSPORT=$((BASEPORT + 2))
DBPORT=$((BASEPORT + 3))
ROOT=$GOPATH/src/github.schibsted.io/Yapo/goms
BASEURL=localhost:$MSPORT
USER=$(whoami)
DB="psql -h localhost -p $DBPORT goms-db"

postgres_test_snap_file=$ROOT/database/snap-tests.sql


if [ -z "$GOPATH" ]
then
	echo "variable GOPATH is not defined!"
	exit -1
fi

case "$1" in
	setup)
		echo ">> $DBPORT"
		echo ">> $MSPORT"
		echo ">> $ROOT"
		echo ">> $BASEURL"
		echo ">> $USER"
		echo
	;;
	*)
		prog=$(basename $0)
		echo $"Usage: $prog {setup}"
		RETVAL=2
esac
