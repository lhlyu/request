package request

import (
	"context"
	"net"
	"net/http"
	"time"
)

func test() {
	r := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, time.Second)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig:       nil,
			DisableKeepAlives:     nil,
		},
	}
	r.Get("")
}
