// The MIT License (MIT)
//
// Copyright (c) 2014 Theeraphol Wattanavekin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package request

import (
	"net/http"
)

// go1.13+ 取消了选择性编译
// 浅拷贝一个transport
func (s *SuperAgent) safeModifyTransport() {
	if !s.isClone {
		return
	}
	oldTransport := s.Transport
	s.Transport = &http.Transport{
		Proxy:                  oldTransport.Proxy,
		DialContext:            oldTransport.DialContext,
		DialTLS:                oldTransport.DialTLS,
		TLSClientConfig:        oldTransport.TLSClientConfig,
		TLSHandshakeTimeout:    oldTransport.TLSHandshakeTimeout,
		DisableKeepAlives:      oldTransport.DisableKeepAlives,
		DisableCompression:     oldTransport.DisableCompression,
		MaxIdleConns:           oldTransport.MaxIdleConns,
		MaxIdleConnsPerHost:    oldTransport.MaxIdleConnsPerHost,
		MaxConnsPerHost:        oldTransport.MaxConnsPerHost,
		IdleConnTimeout:        oldTransport.IdleConnTimeout,
		ResponseHeaderTimeout:  oldTransport.ResponseHeaderTimeout,
		ExpectContinueTimeout:  oldTransport.ExpectContinueTimeout,
		TLSNextProto:           oldTransport.TLSNextProto,
		MaxResponseHeaderBytes: oldTransport.MaxResponseHeaderBytes,
		ProxyConnectHeader:     oldTransport.ProxyConnectHeader,
	}
}
