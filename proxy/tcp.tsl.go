package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Addr struct {
	Ip   string
	Port int
}

type Bridge struct {
	Source             net.Conn       `json:"sourceConn"`
	Destination        net.Conn       `json:"destinationConn"`
}

var listeners = make([]net.Listener, 0)
var ConnMap = make(map[net.Conn]*Bridge)
var Lock = &sync.Mutex{}

var IsStart = false

func Start() {
	for _, p := range ProxyConfig {
		go bindProxy(p)
	}
	IsStart = true
	logrus.Infof("监听已启动")
}

func Stop() {
	for k, v := range ConnMap {
		v.Source.Close()
		v.Destination.Close()
		delete(ConnMap, k)
	}
	for _, l := range listeners {
		l.Close()
	}
	listeners = make([]net.Listener, 0)
	IsStart = false
	logrus.Infof("监听已关闭")
}

func bindProxy(p *ProxyConf) {
	bindPort, isTls, tlsConf := getBindPortAndProxypassPort(p)
	if (PortIsOpen(fmt.Sprintf("0.0.0.0:%v",bindPort),3)){
		logrus.Errorf("端口[%v]已经被占用",bindPort)
		return
	}
	var listener net.Listener
	var err error
	if !isTls {
		listener, err = net.ListenTCP(
			"tcp4",
			&net.TCPAddr{
				Port: bindPort,
			},
		)
	} else { //加密的
		cert, err := tls.LoadX509KeyPair(tlsConf.CrtPath, tlsConf.KeyPath)
		if err != nil {
			log.Println(err)
			return
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", fmt.Sprintf(":%v", bindPort), config)
		if err != nil {
			log.Println(err)
			return
		}
	}
	listeners = append(listeners, listener)
	if !checkError(err) {
		return
	}
	for {
		conn, er := listener.Accept()
		if er != nil {
			break
		}
		go handle(conn, p)
	}
}

/*
* 处理转发
* proxyPort 需要转发的目的端口
*/
func handle(sourceConn net.Conn, p *ProxyConf) {
	defer Lock.Unlock()
	Lock.Lock()
	destination := getDestination(p)
	tcpAddr_dest, err := net.ResolveTCPAddr("tcp4", destination)
	destConn, err := net.DialTCP("tcp", nil, tcpAddr_dest)
	if err != nil {
		logrus.Errorf("DialTCP [%v] error => %v", destination, err)
		return
	}

	bridge := &Bridge{
		Source: sourceConn,
		Destination: destConn,
	}
	ConnMap[sourceConn] = bridge

	go func() {
		defer ReleaseConn(sourceConn, destConn)
		buf := make([]byte, 8)
		io.CopyBuffer(sourceConn, destConn, buf)
	}()

	go func() {
		defer ReleaseConn(sourceConn,destConn)
		buf := make([]byte, 8)
		io.CopyBuffer(destConn, sourceConn, buf)
	}()
}

func checkError(err error) bool {
	if err != nil {
		//logrus.Error("checkError Fatal error: %v", err)
		return false
	}
	return true
}

/* 获取监听端口和转发端口
* bindPort  监听端口，也是源端口
* proxyPort 需要转发的目的端口
*/

func getBindPortAndProxypassPort(p *ProxyConf) (bindPort int, isTls bool, tlsConf *TlsConf) {
	bindPort, _ = strconv.Atoi(p.Source)
	isTls = p.Tls
	tlsConf = p.TlsCf
	return
}

var curDestIndex = 0;
var destCount = 0;
func getDestination(p *ProxyConf) (destination string) {
	destinations := p.Destinations
	if curDestIndex >= destCount {
		curDestIndex = 0
	}
	if strings.Contains(destinations,","){
		destSlice := strings.Split(destinations,",")
		for i,_ := range destSlice {
			d := destSlice[i]
			if PortIsOpen(d,10) {
				i += 1;
				destination = d
				break
			}
		}
	}
	return
}

func ReleaseConn(sourceConn, destConn net.Conn) {
	defer Lock.Unlock()
	Lock.Lock()
	if ConnMap[sourceConn] != nil {
		sourceConn.Close()
		destConn.Close()
		delete(ConnMap, sourceConn)
		logrus.Infof("release connect => %v", sourceConn.RemoteAddr().String())
	}
}
