package main

import (
	"tcpProxy/proxy"
	"tcpProxy/tool"
	"tcpProxy/tpl"
)

func main() {
	tool.InitYmlConfig()
	proxy.LoadFromConfig()
	proxy.Start()
	tpl.ApiStart()
}
