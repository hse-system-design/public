#!/usr/bin/env bash

set -ex

docker exec -it influxdb influx setup
