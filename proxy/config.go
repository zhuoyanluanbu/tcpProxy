package proxy

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type ProxyConf struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Tls         bool    `json:"tls"`
	TlsCf       *TlsConf `json:"tlsCf"`
}
type TlsConf struct {
	CrtPath string `json:"crtPath"`
	KeyPath string	`json:"keyPath"`
}

var configFilePath = "./config.json"
var ProxyConfig = make([]*ProxyConf,0)

func SaveToConfig() {
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		logrus.Error(err.Error())
	}

	b,_ := json.Marshal(&ProxyConfig)
	_, err = f.Write([]byte(""))
	_, err = f.Write(b)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func LoadFromConfig() {
	f, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		logrus.Error(err.Error())
	}
	b, err := ioutil.ReadFile(configFilePath) // just pass the file name
	if err != nil {
		logrus.Error(err)
	}
	json.Unmarshal(b,&ProxyConfig)
}
