package tool

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

var configMap = make(map[string]interface{})

func InitYmlConfig() error{
	path := "app.yml"
	ymlFile,err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	yaml.Unmarshal(ymlFile,configMap);
	if len(configMap) <= 0 {
		return errors.New("cannot load " + path)
	}
	return nil
}

func GetString(key,defaultt string) (result string) {
	result = defaultt
	if !strings.Contains(key,"app.") {
		key = "app."+key
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("GetString err => %s\n", err)
		}
	}()
	keysPath := strings.Split(key,".")
	var res interface{} = nil
	for i,k := range keysPath{
		if i == 0 {
			res = configMap[k].(map[interface{}]interface{})
		}else {
			res = res.(map[interface{}]interface{})[k]
		}

	}
	logrus.Infof("%v -> %v",key,res)
	if res == nil {
		return
	}
	result = fmt.Sprintf("%v",res)
	return
}

func GetInt(key string,defaultt int) int {
	if s:=GetString(key,""); s != "" {
		if n,err := strconv.Atoi(s);err == nil {
			return n;
		}
		return defaultt
	}
	return defaultt
}

func GetStrSliceR(key string,defaultt []string) (result []string) {
	result = defaultt
	if !strings.Contains(key,"app.") {
		key = "app."+key
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("GetString err => %s\n", err)
		}
	}()
	keysPath := strings.Split(key,".")
	var res interface{} = nil
	for i,k := range keysPath{
		if i == 0 {
			res = configMap[k].(map[interface{}]interface{})
		}else {
			res = res.(map[interface{}]interface{})[k]
		}

	}
	logrus.Infof("%v -> %v",key,res)
	if res == nil {
		return
	}
	sli,ok := res.([]interface{})
	if ok {
		for _,s := range sli {
			result = append(result,s.(string))
		}
	}
	return
}