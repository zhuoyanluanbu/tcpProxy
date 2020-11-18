#!/bin/bash

docker build -t hyc/tcpproxy:latest .
mkdir ./bin
docker save -o ./bin/tcpproxy.tar hyc/tcpproxy:latest