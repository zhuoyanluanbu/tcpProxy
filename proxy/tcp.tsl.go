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
	Source             *Addr
	Destination        *Addr
	OuterSelfConn      net.Conn
	SelfDestConn       net.Conn
	WriteToDestination chan []byte
	WriteToSource      chan []byte
}

var listeners = make([]net.Listener, 0)
var connMap = make(map[net.Conn]*Bridge)
var lock = &sync.Mutex{}

var IsStart = false

func Start() {
	for _, p := range ProxyConfig {
		go bindProxy(p)
	}
	IsStart = true
	logrus.Infof("监听已启动")
}

func Stop() {
	for k, v := range connMap {
		v.SelfDestConn.Close()
		v.OuterSelfConn.Close()
		delete(connMap, k)
	}
	for _, l := range listeners {
		l.Close()
	}
	listeners = make([]net.Listener, 0)
	IsStart = false
	logrus.Infof("监听已关闭")
}

func bindProxy(p *ProxyConf) {
	bindPort, destination, isTls, tlsConf := getBindPortAndProxypassPort(p)
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
			continue
		}
		go handle(conn, bindPort, destination)
	}
}

/*
* 处理转发
* proxyPort 需要转发的目的端口
*/
func handle(sourceConn net.Conn, bindPort int, destination string) {
	defer lock.Unlock()
	lock.Lock()

	destination_ip_port := strings.Split(destination, ":")
	destination_ip := destination_ip_port[0]
	destination_port, _ := strconv.Atoi(destination_ip_port[1])

	//destConn := toDestination(sourceConn, destination)
	tcpAddr_dest, err := net.ResolveTCPAddr("tcp4", destination)
	destConn, err := net.DialTCP("tcp", nil, tcpAddr_dest)
	if err != nil {
		logrus.Errorf("DialTCP [%v] error => %v", destination, err)
		return
	}

	bridge := &Bridge{
		OuterSelfConn: sourceConn,
		SelfDestConn:  destConn,
		Source: &Addr{
			Ip:   sourceConn.LocalAddr().String(),
			Port: bindPort,
		},
		Destination: &Addr{
			Ip:   destination_ip,
			Port: destination_port,
		},
		WriteToDestination: make(chan []byte, 1),
		WriteToSource:      make(chan []byte, 1),
	}
	connMap[sourceConn] = bridge

	go func() {
		buf := make([]byte, 2048)
		wr, er := io.CopyBuffer(sourceConn, destConn, buf)
		logrus.Infof("toDestConn WR,Err => %v, %v", wr, er)
		if er != nil {
			releaseConn(sourceConn,destConn)
		}
	}()

	go func() {
		buf := make([]byte, 2048)
		wr, er := io.CopyBuffer(destConn, sourceConn, buf)
		logrus.Infof("toSourceConn WR,Err => %v, %v", wr, er)
		if er != nil {
			releaseConn(sourceConn,destConn)
		}
	}()
}


func checkError(err error) bool {
	if err != nil {
		logrus.Error("checkError Fatal error: %s", err)
		return false
	}
	return true
}

/* 获取监听端口和转发端口
* bindPort  监听端口，也是源端口
* proxyPort 需要转发的目的端口
*/
func getBindPortAndProxypassPort(p *ProxyConf) (bindPort int, destination string, isTls bool, tlsConf *TlsConf) {
	bindPort, _ = strconv.Atoi(p.Source)
	destination = p.Destination
	isTls = p.Tls
	tlsConf = p.TlsCf
	return
}

func releaseConn(sourceConn,destConn net.Conn)  {
	defer lock.Unlock()
	lock.Lock()
	sourceConn.Close()
	destConn.Close()
	delete(connMap,sourceConn)
}
