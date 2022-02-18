## networking
Address range to be used by private networks are:

* 10.0.0.0 to 10.255.255.255
* 172.16.0.0 to 172.31.255.255
* 192.168.0.0 to 192.168.255.255

docker network create --driver=bridge --subnet=10.0.0.0/24 testnet will use 10.0.0.[1-254]

I think we can get away with using /27 by default (that's 30 usable addresses) and make the medic throw an error if more than 30 containers are required by the service.

It can be changed in the config anyway (TODO: maybe configurable subnet per service?) and not just globally.

# Nest

[![Tests](https://github.com/redwebcreation/nest/actions/workflows/tests.yml/badge.svg?branch=next)](https://github.com/redwebcreation/nest/actions/workflows/tests.yml)
[![Static](https://github.com/redwebcreation/nest/actions/workflows/static.yml/badge.svg)](https://github.com/redwebcreation/nest/actions/workflows/static.yml)
[![CodeBeat badge](https://codebeat.co/badges/7171e9ea-53d7-4c81-82bf-a9a2f222b027)](https://codebeat.co/projects/github-com-redwebcreation-nest-next)
[![Go Report Card](https://goreportcard.com/badge/github.com/redwebcreation/nest)](https://goreportcard.com/report/github.com/redwebcreation/nest)
[![codecov](https://codecov.io/gh/redwebcreation/nest/branch/next/graph/badge.svg?token=DWSP4O0YO8)](https://codecov.io/gh/redwebcreation/nest)
![PRs not welcome](https://img.shields.io/badge/PRs-not%20welcome-red)

#### Documentation Status

The goal is to write a lot and then eventually make it more concise and improve upon it.

**VERY MUCH WIP, JUST RANDOM THINGS**

## Requirements

* docker
* git

## What is Nest?

Nest is a tool to help you manage many applications (called "services" from now on) on a single server. You can think of
it as a supercharged docker-compose.

Features:

* zero downtime deployments
* built-in reverse proxy
* versioned configuration
* powerful configuration diagnosis (if anything looks wrong in your configuration, nest will SCREAM LOUDLY)
* an api to trigger deployments automatically (CD [What's Continous Deployment (link needed)]() with a single api call)

## Why use Nest?

Nest is the perfect middle ground between messy configuration files all over your server and kubernetes.

## When not to use Nest?

* You have more than two servers

  If you have exactly two servers, you can still use nest very effectively and make your architecture redundant by
  running them in a Active-Active configuration (or Active-Passive if one is less powerful)
  . [(link needed)]()
