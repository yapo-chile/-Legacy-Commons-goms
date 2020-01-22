#!/usr/bin/env bash

echoTitle "Stopping Docker containers"
docker-compose -f docker/docker-compose.yml -p ${APPNAME} down
