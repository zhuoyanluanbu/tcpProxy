package tpl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"tcpProxy/proxy"
)

type ConnDisplay struct {
	Name   string `json:"name"`
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

var httpServer = http.Server{
	Addr: "127.0.0.1:8081",
}

func TplStart() {
	var index = func (w http.ResponseWriter, r *http.Request) {
		t1, err := template.ParseFiles("./html/index.html")
		if err != nil {
			panic(err)
		}
		t1.Execute(w, "")
	}
	http.HandleFunc("/index", index)
}

func ApiStart() {

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to Golang"))
	})

	//获取链接
	http.HandleFunc("/getAliveConns", func(writer http.ResponseWriter, request *http.Request) {
		connDisplays := make([]*ConnDisplay, 0)
		for k, v := range proxy.ConnMap {
			mem_addr := fmt.Sprintf("%v", &k)
			cd := &ConnDisplay{
				Name:   mem_addr,
				Remote:	v.Source.RemoteAddr().String(),
				Local: v.Destination.RemoteAddr().String(),
			}
			connDisplays = append(connDisplays, cd)
		}
		jsonB, _ := json.Marshal(connDisplays)
		writer.Write(jsonB)
	})

	TplStart()

	httpServer.ListenAndServe()
}
