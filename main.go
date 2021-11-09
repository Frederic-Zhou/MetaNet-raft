package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"metanet/node"
	"metanet/rpc"

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
	rpc.RegisterNodeServer(s, &n)
	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}

}

func CreateConn(address string, timeout time.Duration) (nodeClient *rpc.NodeClient, conn *grpc.ClientConn, err error) {

	//建立链接
	conn, err = grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)

	return &nodeclient, conn, err

	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

	// appendEntries, err := nodeClient.AppendEntries(ctx, &rpc.EntriesArguments{})
	// if err != nil {
	// 	log.Printf("user index could not greet: %v", err)
	// }

	// fmt.Println(appendEntries)
}
