// MIT License
//
// Copyright (c) 2019 Stefan Wichmann
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package lanscan contains a blazing fast port scanner for local networks
package network

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func GetGrpcClientIP(ctx context.Context) (fromip string, selfip string) {

	if pr, ok := peer.FromContext(ctx); ok {

		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			fromip = tcpAddr.IP.String()
		} else {
			fromip = pr.Addr.String()
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		if clientID, ok := md["clientid"]; ok {
			fromip = clientID[0]

		}

		if authoritys, ok := md[":authority"]; ok {
			selfip = strings.Split(authoritys[0], ":")[0]

		}
	}

	return
}
