# Minecraft exporter

[![Test & Build](https://github.com/rebelcore/minecraft_exporter/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/rebelcore/minecraft_exporter/actions/workflows/test.yml)
[![Current Release](https://img.shields.io/github/v/release/rebelcore/minecraft_exporter)](https://github.com/rebelcore/minecraft_exporter/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/rebelcore/minecraft-exporter)](https://hub.docker.com/r/rebelcore/minecraft-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/rebelcore/minecraft_exporter)](https://goreportcard.com/report/github.com/rebelcore/minecraft_exporter)

Prometheus exporter for Minecraft metrics exposed
in Go with pluggable metric collectors.


## Installation and Usage

If you are new to Prometheus and `minecraft_exporter` there is
a [simple step-by-step guide](https://docs.rebelcore.org/guides/minecraft/exporter).

The `minecraft_exporter` listens on HTTP port 9940 by default.
See the `--help` output for more options.


### Ansible

Coming Soon!


### Docker

The `minecraft_exporter` is designed to monitor your Minecraft players real time locations.

For situations where containerized deployment is needed, you will
need to set the RCON address flag to use the docker container hostname.

```bash
docker run -d \
  -p 9940:9940 \
  rebelcore/minecraft-exporter:latest \
  --rcon.address=minecraft:25575
```

For Docker compose, similar flag changes are needed.

```yaml
---
services:
  minecraft_exporter:
    image: rebelcore/minecraft-exporter:latest
    container_name: minecraft_exporter
    command:
      - '--rcon.address=minecraft:25575'
    ports:
      - 9940:9940
    restart: unless-stopped
```


## Collectors

There is varying support for collectors.
The tables below list all existing collectors.

Collectors are enabled by providing a `--collector.<name>` flag.
Collectors that are enabled by default can be disabled
by providing a `--no-collector.<name>` flag.
To enable only some specific collector(s),
use `--collector.disable-defaults --collector.<name> ...`.


### Enabled by default

| Name   | Description                                                           |
|--------|-----------------------------------------------------------------------|
| player | Exposes player username, dimension, XP, and coordinates in real time. |
| system | Exposes if the Minecraft server is online or not.                     |


### Filtering enabled collectors

The `minecraft_exporter` will expose all metrics from enabled collectors
by default. This is the recommended way to collect metrics to avoid errors.

For advanced use the `minecraft_exporter` can be passed an optional list
of collectors to filter metrics. The `collect[]` parameter may be used
multiple times. In Prometheus configuration you can use this syntax under
the [scrape config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#<scrape_config>).

```
  params:
    collect[]:
      - foo
      - bar
```

This can be useful for having different Prometheus servers collect
specific metrics from nodes.


## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* RHEL/CentOS: `glibc-static` package.

Building:

    git clone https://github.com/rebelcore/minecraft_exporter.git
    cd minecraft_exporter
    make build
    ./minecraft_exporter <flags>

To see all available configuration flags:

    ./minecraft_exporter --help


## Running tests

    make test


## TLS endpoint

**EXPERIMENTAL**

The exporter supports TLS via a new web configuration file.

```console
./minecraft_exporter --web.config.file=web-config.yml
```

See the [exporter-toolkit web-configuration](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md) for more details.
