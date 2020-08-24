package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"strconv"
	"strings"
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

var listeners = make([]net.Listener,0)
var connMap = make(map[net.Conn]*Bridge)

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
		delete(connMap,k)
	}
	for _,l := range listeners {
		l.Close()
	}
	listeners = make([]net.Listener,0)
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
	}else { //加密的
		cert, err := tls.LoadX509KeyPair(tlsConf.CrtPath, tlsConf.KeyPath)
		if err != nil {
			log.Println(err)
			return
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", fmt.Sprintf(":%v",bindPort), config)
		if err != nil {
			log.Println(err)
			return
		}
	}
	listeners = append(listeners, listener)
	checkError(err)
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
	destination_ip_port := strings.Split(destination, ":")
	destination_ip := destination_ip_port[0]
	destination_port, _ := strconv.Atoi(destination_ip_port[1])

	destConn := toDestination(sourceConn, destination)

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

	sendToDest(sourceConn)
	sendToSource(sourceConn)

	for {
		bsLen := 1024 * 16
		bs := make([]byte, bsLen)
		_, err := sourceConn.Read(bs)
		if err != nil {
			logrus.Errorf("read from source[%v]=> %v", sourceConn.LocalAddr().String()+":"+strconv.Itoa(bindPort), err)
			delete(connMap, sourceConn)
			sourceConn.Close()
			break
		}
		bridge.WriteToDestination <- bs
	}
}

/*
* 转发到目的
*/
func toDestination(sourceConn net.Conn, destination string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", destination)
	if err != nil {
		logrus.Errorf("ResolveTCPAddr [%v] error => %v", destination, err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logrus.Errorf("DialTCP [%v] error => %v", destination, err)
	}
	go func() {
		for {
			bsLen := 1024 * 16
			bs := make([]byte, bsLen)
			_, err := conn.Read(bs)
			if err != nil {
				logrus.Errorf("read from destination[%v] err => %v", destination, err)
				break
			}
			bri := connMap[sourceConn]
			if bri != nil {
				bri.WriteToSource <- bs
			} else {
				break
			}
		}
	}()
	logrus.Infof("connect [%v] success", destination)
	return conn
}

func sendToDest(sourceConn net.Conn) {
	go func() {
		bri := connMap[sourceConn]
		if bri != nil {
			for {
				bs := <-bri.WriteToDestination
				conn := bri.SelfDestConn
				conn.Write(bs)
			}
		}
	}()
}

func sendToSource(sourceConn net.Conn) {
	go func() {
		bri := connMap[sourceConn]
		if bri != nil {
			for {
				bs := <-bri.WriteToSource
				conn := bri.OuterSelfConn
				conn.Write(bs)
			}
		}
	}()
}

func checkError(err error) {
	if err != nil {
		logrus.Fatalf("checkError Fatal error: %s", err)
	}
}

/* 获取监听端口和转发端口
* bindPort  监听端口，也是源端口
* proxyPort 需要转发的目的端口
*/
func getBindPortAndProxypassPort(p *ProxyConf) (bindPort int, destination string,isTls bool,tlsConf *TlsConf) {
	bindPort, _ = strconv.Atoi(p.Source)
	destination = p.Destination
	isTls = p.Tls
	tlsConf = p.TlsCf
	return
}
