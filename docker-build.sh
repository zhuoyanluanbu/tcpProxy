#!/bin/bash

./build-all.sh
docker build -t cetciot/mqtt-broker-proxy:latest .
mkdir ./bin
docker save -o ./bin/mqtt-broker-proxy.tar cetciot/mqtt-broker-proxy:latest