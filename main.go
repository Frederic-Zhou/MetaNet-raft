package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"metanet/node"
	"metanet/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	n := node.Node{}
	n.Timeout = node.MinTimeout + time.Duration(rand.Intn(node.MaxTimeout-node.MinTimeout))

	// 监听本地端口
	lis, err := net.Listen("tcp", ":8800")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}

	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	rpc.RegisterNodeServer(s, &n)
	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}

}
