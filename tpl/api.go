package tpl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"tcpProxy/proxy"
)

type ConnDisplay struct {
	ID   string `json:"id"`
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
				ID:   mem_addr,
				Remote:	v.Source.RemoteAddr().String(),
				Local: v.Destination.RemoteAddr().String(),
			}
			connDisplays = append(connDisplays, cd)
		}
		jsonB, _ := json.Marshal(connDisplays)
		writer.Write(jsonB)
	})

	//获取证书文件夹下面的文件
	http.HandleFunc("/getCertFileNames", func(writer http.ResponseWriter, request *http.Request) {
		names := getFileNames("./certs")
		jsonB, _ := json.Marshal(names)
		writer.Write(jsonB)
	})

	//上传证书
	http.HandleFunc("/uploadCert", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		fileKey := "file"
		filename := request.Header.Get("filename")

		for _,n := range getFileNames("./certs") {
			if n == filename {
				os.Remove(n)
			}
		}
		fileDest,_ := os.OpenFile(filename,os.O_RDWR|os.O_CREATE,777)
		fileSource,_,err := request.FormFile(fileKey)
		defer fileSource.Close()
		if err != nil {
			return
		}
		if _, err := io.Copy(fileDest, fileSource); err != nil {
			return
		}
		writer.Write([]byte("ok"))
	})



	TplStart()
	httpServer.ListenAndServe()
}

func getFileNames(path string) []string  {
	names := []string{}
	files,_ := ioutil.ReadDir(path)
	for _,f := range files {
		names = append(names, f.Name())
	}
	return names
}