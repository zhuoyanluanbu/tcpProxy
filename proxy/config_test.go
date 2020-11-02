package proxy

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSaveToConfig(t *testing.T) {
	ProxyConfig = append(
		ProxyConfig,
		&ProxyConf{
			Source:      "127.0.0.1:2020",
			Destinations: "127.0.0.1:2200",
			Tls:         false,
		},)
	SaveToConfig()
}

func TestLoadFromConfig(t *testing.T) {
	LoadFromConfig()
	b,_:=json.Marshal(ProxyConfig)
	fmt.Printf("TestLoadFromConfig => %v",string(b))
}
