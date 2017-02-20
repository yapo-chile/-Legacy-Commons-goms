#|/bin/bash

. scripts/functions

if [ -z "$CWD" ]; then
    CWD=`pwd`
else
    CWD=$CWD
fi

PSQL=/usr/pgsql-9.3/bin/pg_ctl
INITDB=/usr/pgsql-9.3/bin/initdb
CREATEDB=/usr/pgsql-9.3/bin/createdb
DROPDB=/usr/pgsql-9.3/bin/dropdb
DATAPATH=$CWD/database/data

PGSQL_PORT=$2
PGSQL_DATABASE=$3
OPTIONS="-s -D $DATAPATH/ -l $CWD/logs/postgres.log -o \"-p $PGSQL_PORT\" "

start() {
	if [ ! -f $DATAPATH/PG_VERSION ]; then
		CALL_INIT=1
		initDB
	fi
	echo -n "Starting postgresql: "
	daemon $PSQL $OPTIONS "start"
	echo

	sleep 5
	if [ "$CALL_INIT" ]; then
		createDB
	fi
}

initDB() {
	mkdir -p $DATAPATH
	rm -rf $DATAPATH/*
	$INITDB $DATAPATH
}

createDB() {
	$CREATEDB -h localhost -p $PGSQL_PORT $PGSQL_DATABASE
	RETVAL=$?
	psql -h localhost -p $PGSQL_PORT $PGSQL_DATABASE < $CWD/database/structure.sql
	echo
}

clean() {
	$PSQL status -D $DATAPATH > /dev/null 2>&1
	OUTPUT=$?
	if [ "$OUTPUT" -eq 0 ]; then
		echo "ERROR: DATABASE RUNNING"
		return -1
	else
		rm -Rf "$DATAPATH"
	fi
}

status() {
	$PSQL status -D $DATAPATH > /dev/null 2>&1
	OUTPUT=$?
	if [ "$OUTPUT" -eq 0 ]; then
		echo "DATABASE RUNNING"
	else
		echo "DATABASE NOT RUNNING"
	fi
}

stop() {
	echo -n "Stoping postgresql: "
	daemon $PSQL $OPTIONS "stop"
	echo
}

case "$1" in
	start)
		start
		;;
	stop)
		stop
		;;
	status)
		status
		;;
	clean)
		clean
		;;
esac
