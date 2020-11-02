package proxy

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func PortIsOpen(addr string, timeout int) bool {
	con, err := net.DialTimeout("tcp", addr, time.Second * time.Duration(timeout))
	if err != nil {
		return false
	}
	con.Close()
	return true
}

func Telnet(action []string, ip string, timeout int) (buf []byte, err error) {
	con, err := net.DialTimeout("tcp", ip, time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(time.Second * 5))
	for _, v := range action {
		l := strings.SplitN(v, "_", 2)
		if len(l) < 2 {
			return
		}
		switch l[0] {
		case "r":
			var n int
			n, err = strconv.Atoi(l[1])
			if err != nil {
				return
			}
			p := make([]byte, n)
			n, err = con.Read(p)
			if err != nil {
				return
			}
			buf = append(buf, p[:n]...)
			fmt.Println(buf)
		case "w":
			_, err = con.Write([]byte(l[1]))
		}
	}
	return
}