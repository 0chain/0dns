#!/bin/sh

set -e

PWD=`pwd`

echo Starting 0dns ...

docker-compose -p 0dns -f ./docker.local/docker-compose-no-daemon.yml up
