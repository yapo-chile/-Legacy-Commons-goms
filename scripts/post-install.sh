#!/bin/bash

HOSTNAME=$(hostname --short)
FULLHOSTNAME=$(hostname)
FILE='/opt/__API_NAME__/conf/conf.json'

sed -i "s/__SERVER_NAME__/$FULLHOSTNAME/g" $FILE
sed -i "s/__SERVICE_PORT__/6969/g" $FILE
sed -i "s/__DB_NAME__/goms-db/g" $FILE
sed -i "s/__DB_SERVER_PORT__/5432/g" $FILE
sed -i "s/__DB_USER__/gomsms/g" $FILE
sed -i "s/__DB_PASSWD__/XgG5M_Qe3/g" $FILE

sed -i "s/__CACHE_ENABLED__/true/g" $FILE
sed -i "s/__CACHE_PASSWORD__//g" $FILE
sed -i "s/__CACHE_PORT__/6399/g" $FILE
sed -i "s/__CACHE_DATABASE_NUMBER__/0/g" $FILE
sed -i "s/__CACHE_MAX_IDLE__/10/g" $FILE
sed -i "s/__CACHE_TIMEOUT__/60/g" $FILE
sed -i "s/__CACHE_MAX_ACTIVE__/100/g" $FILE
sed -i "s/__CACHE_WAIT__/false/g" $FILE

sed -i "s/__SYSLOG_ENABLED__/true/g" $FILE
sed -i "s/__SYSLOG_IDENTITY__/goms/g" $FILE
sed -i "s/__STDLOG_ENABLED__/false/g" $FILE
sed -i "s/__LOG_LEVEL__/1/g" $FILE

case "$HOSTNAME" in
	"lvdev-payment01") # poya1
		sed -i "s/__DB_HOSTNAME__/lvdev-db01.schibsted.cl/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav01/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $FILE
	;;
	"lvdev-payment02") # poya2
		sed -i "s/__DB_HOSTNAME__/lvdev-db02.schibsted.cl/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav02/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $FILE
	;;
	"lvdev-payment03") # poya3
		sed -i "s/__DB_HOSTNAME__/lvdev-db03.schibsted.cl/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav03/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $FILE
	;;
	"lvdev-payment04") # poya4
		sed -i "s/__DB_HOSTNAME__/lvdev-db04.schibsted.cl/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/lvdev-dav04/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/600/g" $FILE
	;;
	"ch32stg") # STG
		sed -i "s/__DB_HOSTNAME__/ch6stg/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/ch36stg/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/3600/g" $FILE
	;;
	"ch42") # prod
		sed -i "s/__DB_HOSTNAME__/ch6/g" $FILE
		sed -i "s/__CACHE_DB_HOSTNAME__/ch36/g" $FILE
		sed -i "s/__CACHE_EXPIRATION__/3600/g" $FILE
	;;
esac
