# 0dns - A DNS service for Züs clients

This service is responsible for connecting to the network and fetching all the magic blocks from the network which are saved in the DB.
There is a [Network API](#network-api) which can be used to get latest set of miners and sharders for the network using the magic block.

## Table of Contents
- [Züs Overview](#züs-overview)
- [Setup](#setup)
- [Buildding and starting the node](#building-and-starting-the-node)
- [Point to another blockchain](#point-to-another-blockchain)
- [Exposed APIs](#exposed-apis)
- [Config changes](#config-changes)
- [Cleanup](#cleanup)
- [Network issue](#network-issue)

## Züs Overview
[Züs](https://zus.network/) is a high-performance cloud on a fast blockchain offering privacy and configurable uptime. It is an alternative to traditional cloud S3 and has shown better performance on a test network due to its parallel data architecture. The technology uses erasure code to distribute the data between data and parity servers. Züs storage is configurable to provide flexibility for IT managers to design for desired security and uptime, and can design a hybrid or a multi-cloud architecture with a few clicks using [Blimp's](https://blimp.software/) workflow, and can change redundancy and providers on the fly.

For instance, the user can start with 10 data and 5 parity providers and select where they are located globally, and later decide to add a provider on-the-fly to increase resilience, performance, or switch to a lower cost provider.

Users can also add their own servers to the network to operate in a hybrid cloud architecture. Such flexibility allows the user to improve their regulatory, content distribution, and security requirements with a true multi-cloud architecture. Users can also construct a private cloud with all of their own servers rented across the globe to have a better content distribution, highly available network, higher performance, and lower cost.

[The QoS protocol](https://medium.com/0chain/qos-protocol-weekly-debrief-april-12-2023-44524924381f) is time-based where the blockchain challenges a provider on a file that the provider must respond within a certain time based on its size to pass. This forces the provider to have a good server and data center performance to earn rewards and income.

The [privacy protocol](https://zus.network/build) from Züs is unique where a user can easily share their encrypted data with their business partners, friends, and family through a proxy key sharing protocol, where the key is given to the providers, and they re-encrypt the data using the proxy key so that only the recipient can decrypt it with their private key.

Züs has ecosystem apps to encourage traditional storage consumption such as [Blimp](https://blimp.software/), a S3 server and cloud migration platform, and [Vult](https://vult.network/), a personal cloud app to store encrypted data and share privately with friends and family, and [Chalk](https://chalk.software/), a high-performance story-telling storage solution for NFT artists.

Other apps are [Bolt](https://bolt.holdings/), a wallet that is very secure with air-gapped 2FA split-key protocol to prevent hacks from compromising your digital assets, and it enables you to stake and earn from the storage providers; [Atlus](https://atlus.cloud/), a blockchain explorer and [Chimney](https://demo.chimney.software/), which allows anyone to join the network and earn using their server or by just renting one, with no prior knowledge required.

## Setup

Clone the repo and run the following command inside the cloned directory

```
./docker.local/bin/init.sh
```

## Building and Starting the Node

If there is new code, do a git pull and run the following command

```
./docker.local/bin/build.sh
```

Run the container using

```
./docker.local/bin/start.sh
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
./docker.local/bin/clean.sh
```

## Network issue

If there is no test network, run the following command

```
docker network create --driver=bridge --subnet=198.18.0.0/15 --gateway=198.18.0.255 testnet0
```

## Troubleshooting

1. When running locally, if `http://localhost:9091/network` is returning below. Some commands might fail.

```
{
  "miners": [
    "https://198.18.0.71/",
    "https://198.18.0.72/",
    "https://198.18.0.74/",
    "https://198.18.0.73/"
  ],
  "sharders": [
    "https://198.18.0.81/"
  ]
}
```

To fix the miner and sharder URLs so it works locally, update `docker.local/config/0dns.yaml` to disable both `use_https` and `use_path` (set to `false`).
