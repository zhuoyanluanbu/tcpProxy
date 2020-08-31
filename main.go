package main

import (
	"os/exec"
	"tcpProxy/proxy"
	"tcpProxy/tool"
	"tcpProxy/tpl"
)

func main() {
	tool.InitYmlConfig()
	proxy.LoadFromConfig()
	tpl.ApiStart(openBrowser)
}

func openBrowser (url string)  {
	exec.Command(`cmd`, `/c`, `start`, url).Start()
	exec.Command(`xdg-open`, url).Start()
	exec.Command(`open`, url).Start()
}
