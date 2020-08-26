package main

import (
	"tcpProxy/tool"
	"tcpProxy/tpl"
)

func main() {
	tool.InitYmlConfig()
	//proxy.LoadFromConfig()
	//proxy.Start()
	tpl.ApiStart()
}
