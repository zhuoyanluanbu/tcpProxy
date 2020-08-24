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
	//gui.RunWindow()
	tpl.RunTplTest()
	//test.ListenTls()
}
