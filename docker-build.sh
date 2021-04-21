#!/bin/bash

./build-all.sh
docker build -t cetciot/tcp-proxy:latest .
mkdir ./bin
docker save -o ./bin/tcp-proxy.tar cetciot/tcp-proxy:latest