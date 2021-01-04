package tpl

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"tcpProxy/proxy"
	"time"
)

type ConnDisplay struct {
	ID     string `json:"id"`
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

type ProxyConfigVo struct {
	Port1 string `json:"port1"`
	Ip string `json:"ip"`
	Port2 string	`json:"port2"`
	Tls bool	`json:"tls"`
	CrtPath string	`json:"crtPath"`
	KeyPath string	`json:"keyPath"`
}

var httpServer = http.Server{
	Addr: ":18081",
}

func TplStart() {
	http.Handle("/static/",
		http.StripPrefix("/static/",http.FileServer(http.Dir("./html/static")) ),
	)
	var index = func(w http.ResponseWriter, r *http.Request) {
		t1, err := template.ParseFiles("./html/index.html")
		if err != nil {
			panic(err)
		}
		t1.Execute(w, "")
	}
	http.HandleFunc("/index", index)
}

func ApiStart(complete func(string)) {

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to Golang"))
	})

	http.HandleFunc("/runState",func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("%v",proxy.IsStart)))
	})

	http.HandleFunc("/startup",func(writer http.ResponseWriter, request *http.Request) {
		bf, _ := ioutil.ReadAll(request.Body)
		pcvs := []*ProxyConfigVo{}
		json.Unmarshal(bf,&pcvs)
		proxy.ClearProxyConfig()
		for _,po := range pcvs {
			pc := &proxy.ProxyConf{
				Source: fmt.Sprintf("%v",po.Port1),
				//Destinations: fmt.Sprintf("%v:%v",po.Ip,po.Port2),
				Destinations: fmt.Sprintf("%v",po.Ip),
				Tls: po.Tls,
			}
			if po.Tls {
				pc.TlsCf = &proxy.TlsConf{
					KeyPath: po.KeyPath,
					CrtPath: po.CrtPath,
				}
				if !strings.Contains(po.KeyPath,"./certs") {
					pc.TlsCf.KeyPath = fmt.Sprintf("./certs/%v",po.KeyPath)

				}
				if !strings.Contains(po.KeyPath,"./certs") {
					pc.TlsCf.CrtPath = fmt.Sprintf("./certs/%v",po.CrtPath)
				}

			}
			if po.KeyPath=="" || po.CrtPath == "" {
				pc.Tls = false
			}
			proxy.ProxyConfig = append(proxy.ProxyConfig, pc)
		}
		buff := ""
		if !proxy.IsStart {
			proxy.SaveToConfig()
			for _,p := range proxy.ProxyConfig {
				if proxy.PortIsOpen(fmt.Sprintf("0.0.0.0:%v",p.Source),3) {
					buff += fmt.Sprintf("端口[%v]被占用",p.Source)
				}
			}
		}

		if buff == "" {
			proxy.Start()
			writer.Write([]byte("ok"))
		}else {
			proxy.Stop()
			writer.Write([]byte(buff))
		}

	})

	http.HandleFunc("/shutdown",func(writer http.ResponseWriter, request *http.Request) {
		if proxy.IsStart {
			proxy.Stop()
		}
		writer.Write([]byte("ok"))
	})

	//获取配置
	http.HandleFunc("/getConfig", func(writer http.ResponseWriter, request *http.Request) {
		jsonB, _ := json.Marshal(proxy.ProxyConfig)
		writer.Write(jsonB)
	})

	//获取链接
	http.HandleFunc("/getAliveConns", func(writer http.ResponseWriter, request *http.Request) {
		connDisplays := make([]*ConnDisplay, 0)
		for k, v := range proxy.ConnMap {
			mem_addr := fmt.Sprintf("%v", *&k)
			mem_addr = strings.ReplaceAll(mem_addr,"&{{","")
			mem_addr = strings.ReplaceAll(mem_addr,"}}","")
			cd := &ConnDisplay{
				ID:   mem_addr,
				Remote:	v.Source.RemoteAddr().String(),
				Local: v.Destination.RemoteAddr().String(),
			}
			connDisplays = append(connDisplays, cd)
		}
		sort.Slice(connDisplays, func(i, j int) bool {
			return connDisplays[i].ID < connDisplays[j].ID
		})
		jsonB, _ := json.Marshal(connDisplays)
		writer.Write(jsonB)
	})

	//踢掉某个链接
	http.HandleFunc("/dropConn", func(writer http.ResponseWriter, request *http.Request) {
		id := request.URL.Query().Get("id")
		func(mem string){
			for k, v := range proxy.ConnMap {
				mem_addr := fmt.Sprintf("%v", *&k)
				mem_addr = strings.ReplaceAll(mem_addr,"&{{","")
				mem_addr = strings.ReplaceAll(mem_addr,"}}","")
				if mem == mem_addr {
					proxy.ReleaseConn(v.Source,v.Destination)
				}
			}
		}(id)
		writer.Write([]byte("ok"))
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

		for _, n := range getFileNames("./certs") {
			if n == filename {
				os.Remove(n)
			}
		}
		fileDest, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 777)
		fileSource, _, err := request.FormFile(fileKey)
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

	go func() {
		time.Sleep(time.Second*2)
		url := fmt.Sprintf("http://127.0.0.1%v/index",httpServer.Addr)
		complete(url)
	}()

	fmt.Printf("http 监听 %v",httpServer.Addr)
	httpServer.ListenAndServe()

}

func getFileNames(path string) []string {
	names := []string{}
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names
}


