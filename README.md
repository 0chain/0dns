# 0dns

This service is responsible for connecting to the network and fetching all the magic blocks from the network which are saved in the DB.
There is a [Network API](#network-api) which can be used to get latest set of miners and sharders for the network using the magic block.

## Table of Contents

- [Setup](#setup)
- [Buildding and starting the node](#building-and-starting-the-node)
- [Point to another blockchain](#point-to-another-blockchain)
- [Exposed APIs](#exposed-apis)
- [Config changes](#config-changes)
- [Cleanup](#cleanup)
- [Network issue](#network-issue)

## Setup

Clone the repo and run the following command inside the cloned directory

```
$ ./docker.local/bin/init.sh
```

## Building and Starting the Node

If there is new code, do a git pull and run the following command

```
$ ./docker.local/bin/build.sh
```

Go to the bin directory (cd docker.local/bin) and run the container using

```
$ ./start.sh
```

## Point to another blockchain

You can point the server to any instance of 0chain blockchain you like, Just go to config (docker.local/config) and update the magic_block.json.

By default they are pointing to local network with 4 miners and 1 sharder

You have to get the magic block from the chain you want to connect and update the `magic_block.json` file.

It will use that file to extract miners and sharders value and use that for further operations.

## Exposed APIs

0dns exposes:

### Network API

```
{BASEURL}/network
```

This API can be used to fetch the latest set of miners and sharder to connect.
GoSDK will also uses this API to fetch miner/sharders value on runtime.

Details: https://0chain.net/page-documentation.html#tag/BlockWorker/paths/~1network/get

### Magic Block API

```
{BASEURL}/magic_block
```

This API can be used to fetch the latest magic block data from network.

Details: https://0chain.net/page-documentation.html#tag/BlockWorker/paths/~1magic_block/get

## Config Changes

You can do other config changes as well in 0dns.yaml file itself, Like

- Mongo DB connection URL, DB name and pool size

```
mongo:
  url: mongodb://mongodb:27017
  db_name: block-recorder
  pool_size: 2
```

- Server port

```
port: 9091
```

- Logging info level

```
logging:
  level: "info"
  console: false # printing log to console is only supported in development mode
```

- Worker related config

```
worker:
  magic_block_worker: 5 # in seconds, Frequency to update magic block
```

## Cleanup

Get rid of old data when the blockchain is restarted or if you point to a different network:

```
$ ./docker.local/bin/clean.sh
```

## Network issue

If there is no test network, run the following command

```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```
