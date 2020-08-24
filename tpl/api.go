package tpl

import "net/http"

var HttpServer = http.Server{
	Addr: "127.0.0.1:8080",
}

func ApiStart() {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to Golang"))
	})
}
