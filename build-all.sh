#!/usr/bin/env bash

GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o TcpProxy.exe
GOOS=windows GOARCH=386 go build -x -v -ldflags "-s -w" -o TcpProxy32.exe
GOOS=linux GOARCH=amd64 go build -o TcpProxy_linux64 main.go
go build -o TcpProxy_macos main.go
tar -czvf tcpProxy.tar.gz config.json html/ certs/ TcpProxy.exe TcpProxy32.exe TcpProxy_linux64 TcpProxy_macos
