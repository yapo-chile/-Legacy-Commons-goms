#!/bin/bash

HOSTNAME=$(hostname --short)
FULLHOSTNAME=$(hostname)
CONF_FILE='/opt/__API_NAME__/conf/conf.json'
INIT_SCRIPT='/etc/init.d/__API_NAME__-api'
SERVICE_PORT=16969

case "${HOSTNAME}" in
	"lvdev-payment01") # poya1
		DB_HOSTNAME="lvdev-db01.schibsted.cl"
		CACHE_DB_HOSTNAME="lvdev-dav01"
		CACHE_EXPIRATION=600
	;;
	"lvdev-payment02") # poya2
		DB_HOSTNAME="lvdev-db02.schibsted.cl"
		CACHE_DB_HOSTNAME="lvdev-dav02"
		CACHE_EXPIRATION=600
	;;
	"lvdev-payment03") # poya3
		DB_HOSTNAME="lvdev-db03.schibsted.cl"
		CACHE_DB_HOSTNAME="lvdev-dav03"
		CACHE_EXPIRATION=600
	;;
	"lvdev-payment04") # poya4
		DB_HOSTNAME="lvdev-db04.schibsted.cl"
		CACHE_DB_HOSTNAME="lvdev-dav04"
		CACHE_EXPIRATION=600
	;;
	"ch32stg") # STG
		DB_HOSTNAME="ch6stg"
		CACHE_DB_HOSTNAME="ch36stg"
		CACHE_EXPIRATION=3600
	;;
	"ch42") # prod
		DB_HOSTNAME="ch6"
		CACHE_DB_HOSTNAME="ch36"
		CACHE_EXPIRATION=3600
	;;
esac

sed -i "s/__SERVICE_PORT__/${SERVICE_PORT}/g" $INIT_SCRIPT
sed -e "s/__SERVER_NAME__/$FULLHOSTNAME/g" \
    -e "s/__SERVICE_PORT__/${SERVICE_PORT}/g" \
    -e "s/__SERVICE_PID__/.pid/g" \
    -e "s/__DB_NAME__/goms-db/g" \
    -e "s/__DB_SERVER_PORT__/5432/g" \
    -e "s/__DB_USER__/gomsms/g" \
    -e "s/__DB_PASSWD__/XgG5M_Qe3/g" \
    -i ${CONF_FILE}

sed -e "s/__CACHE_ENABLED__/true/g" \
    -e "s/__CACHE_PASSWORD__//g" \
    -e "s/__CACHE_PORT__/6399/g" \
    -e "s/__CACHE_DATABASE_NUMBER__/0/g" \
    -e "s/__CACHE_MAX_IDLE__/10/g" \
    -e "s/__CACHE_TIMEOUT__/60/g" \
    -e "s/__CACHE_MAX_ACTIVE__/100/g" \
    -e "s/__CACHE_WAIT__/false/g" \
    -i ${CONF_FILE}

sed -e "s/__SYSLOG_ENABLED__/true/g" \
    -e "s/__SYSLOG_IDENTITY__/goms/g" \
    -e "s/__STDLOG_ENABLED__/false/g" \
    -e "s/__LOG_LEVEL__/1/g" \
    -i ${CONF_FILE}

sed -e "s/__DB_HOSTNAME__/${DB_HOSTNAME}/g" \
    -e "s/__CACHE_DB_HOSTNAME__/${CACHE_DB_HOSTNAME}/g" \
    -e "s/__CACHE_EXPIRATION__/${CACHE_EXPIRATION}/g" \
    -i ${CONF_FILE}
