package main

import (
	"fmt"
	"os/exec"
	"tcpProxy/proxy"
	"tcpProxy/tool"
	"tcpProxy/tpl"
)

func main() {
	tool.InitYmlConfig()
	fmt.Println("init yml ...")
	proxy.LoadFromConfig()
	fmt.Println("init config.json ...")
	tpl.ApiStart(openBrowser)
}

func openBrowser (url string)  {
	exec.Command(`cmd`, `/c`, `start`, url).Start()
	exec.Command(`xdg-open`, url).Start()
	exec.Command(`open`, url).Start()
}

