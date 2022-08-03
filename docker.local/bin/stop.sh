#!/bin/sh
PWD=`pwd`

echo Stopping 0dns ...

docker-compose -p 0dns -f ../docker-compose.yml down
