#!/bin/bash

HOSTNAME=$(hostname --short)
FULLHOSTNAME=$(hostname)
CONF_FILE='/opt/__API_NAME__/conf/conf.json'
INIT_SCRIPT='/etc/init.d/__API_NAME__-api'
SERVICE_PORT=7070

sed -i "s/__SERVICE_PORT__/${SERVICE_PORT}/g" $INIT_SCRIPT
sed -i "s/__SERVER_NAME__/$FULLHOSTNAME/g" $CONF_FILE
sed -i "s/__SERVICE_PORT__/${SERVICE_PORT}/g" $CONF_FILE
sed -i "s/__DB_NAME__/goms-db/g" $CONF_FILE
sed -i "s/__DB_SERVER_PORT__/5432/g" $CONF_FILE
sed -i "s/__DB_USER__/gomsms/g" $CONF_FILE
sed -i "s/__DB_PASSWD__/XgG5M_Qe3/g" $CONF_FILE

sed -i "s/__CACHE_ENABLED__/true/g" $CONF_FILE
sed -i "s/__CACHE_PASSWORD__//g" $CONF_FILE
sed -i "s/__CACHE_PORT__/6399/g" $CONF_FILE
sed -i "s/__CACHE_DATABASE_NUMBER__/0/g" $CONF_FILE
sed -i "s/__CACHE_MAX_IDLE__/10/g" $CONF_FILE
sed -i "s/__CACHE_TIMEOUT__/60/g" $CONF_FILE
sed -i "s/__CACHE_MAX_ACTIVE__/100/g" $CONF_FILE
sed -i "s/__CACHE_WAIT__/false/g" $CONF_FILE

sed -i "s/__SYSLOG_ENABLED__/true/g" $CONF_FILE
sed -i "s/__SYSLOG_IDENTITY__/goms/g" $CONF_FILE
sed -i "s/__STDLOG_ENABLED__/false/g" $CONF_FILE
sed -i "s/__LOG_LEVEL__/1/g" $CONF_FILE

case "$HOSTNAME" in
	"lvdev-payment01") # poya1
		sed -i "s/__DB_HOSTNAME__/lvdev-db01.schibsted.cl/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav01/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $CONF_FILE
	;;
	"lvdev-payment02") # poya2
		sed -i "s/__DB_HOSTNAME__/lvdev-db02.schibsted.cl/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav02/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $CONF_FILE
	;;
	"lvdev-payment03") # poya3
		sed -i "s/__DB_HOSTNAME__/lvdev-db03.schibsted.cl/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav03/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $CONF_FILE
	;;
	"lvdev-payment04") # poya4
		sed -i "s/__DB_HOSTNAME__/lvdev-db04.schibsted.cl/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav04/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $CONF_FILE
	;;
	"ch32stg") # STG
		sed -i "s/__DB_HOSTNAME__/ch6stg/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/ch36stg/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/3600/g" $CONF_FILE
	;;
	"ch42") # prod
		sed -i "s/__DB_HOSTNAME__/ch6/g" $CONF_FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/ch36/g" $CONF_FILE
		sed -i "s/__CACHE_EXPIRATION__/3600/g" $CONF_FILE
	;;
esac
