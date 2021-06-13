#!/bin/sh

#mount volumns

VOLUMES_CONFIG=./config VOLUMES_LOG=./0dns/log VOLUMES_MONGO_DATA=./0dns/mongodata docker-compose -p 0dns -f docker.local/docker-compose.yml build --force-rm
docker.local/bin/sync_clock.sh
