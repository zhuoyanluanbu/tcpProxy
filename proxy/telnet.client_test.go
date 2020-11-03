package proxy

import (
	"reflect"
	"testing"
)

func TestPortIsOpen(t *testing.T) {
	type args struct {
		ip      string
		timeout int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name:"TestPortIsOpen",
			args: args{
				ip:      "0.0.0.0:1081",
				timeout: 10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PortIsOpen(tt.args.ip, tt.args.timeout); got != tt.want {
				t.Errorf("PortIsOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelnet(t *testing.T) {
	type args struct {
		action  []string
		ip      string
		timeout int
	}
	tests := []struct {
		name    string
		args    args
		wantBuf []byte
		wantErr bool
	}{
		{
			name:"TestPortIsOpen",
			args: args{
				ip:      "192.168.0.13:1883",
				timeout: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, err := Telnet(tt.args.action, tt.args.ip, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Telnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuf, tt.wantBuf) {
				t.Errorf("Telnet() = %v, want %v", gotBuf, tt.wantBuf)
			}
		})
	}
}
