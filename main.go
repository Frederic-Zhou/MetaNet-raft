package main

import (
	"fmt"
	"net"

	"metanet/node"
	"metanet/node_rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	n := node.Node{}

	// 监听本地端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}

	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	node_rpc.RegisterNodeServer(s, &n)
	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}

}
