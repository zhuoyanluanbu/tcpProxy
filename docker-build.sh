#!/bin/bash

docker build -t cetciot/mqtt-broker-proxy:latest .
mkdir ./bin
docker save -o ./bin/tcpproxy.tar cetciot/mqtt-broker-proxy:latest