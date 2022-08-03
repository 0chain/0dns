#!/bin/sh

#
# clean up without sudo being a member of the docker group
#

set -e

docker-compose                                                \
    -f ./docker.local/docker-clean/docker-clean-compose.yml   \
    up --build docker-clean
