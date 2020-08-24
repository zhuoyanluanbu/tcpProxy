package test

import (
	"testing"
	"time"
)

func TestTlsClient(t *testing.T) {
	for i:=0;i<10;i++{
		time.Sleep(time.Second*3)
		TlsClient()
	}
}