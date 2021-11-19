package node

import (
	context "context"
	"encoding/json"
	"fmt"
	"log"
	"metanet/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client = Node

func (c *Client) ClientRequestCall(cmd []byte, to string, md map[string]string) (result *rpc.ClientResults, err error) {

	result = &rpc.ClientResults{}

	//链接到节点
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", to, PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}

	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), MinTimeout*time.Millisecond)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(md))

	result, err = nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: cmd})

	return
}

func (c *Client) Command(cmd string, args ...string) (r interface{}, err error) {

	switch cmd {
	case "ll":
		return json.Marshal(c.Log)
	case "lc":
		return json.Marshal(c.NodesConfig)
	case CMD_JOIN:

		if len(args) > 1 {
			return c.ClientRequestCall([]byte(CMD_JOIN), args[1], nil)
		} else {
			err = fmt.Errorf("提供一个地址参数")
			return
		}

	}

	return

}

// // func (c *Client) Join() (leaderID string, fastNodeID string) {
// // 	//var validNetworks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

// // 	//得到所有tcp4的网卡的网络环境的IP列表（理论列表）
// // 	allIPs := []string{}
// // 	for _, current := range network.LinkLocalAddresses("tcp4") {
// // 		allIPs = append(allIPs, network.CalculateSubnetIPs(current)...)
// // 	}

// // 	count := len(allIPs)

// // 	hostchan := make(chan string, count)
// // 	resultchan := make(chan []string, count)
// // 	//创建n个联络员
// // 	n := 100
// // 	for i := 0; i < n; i++ {
// // 		go liaison(hostchan, resultchan)
// // 	}

// // 	//轮训所有可连接地址
// // 	for _, host := range allIPs {
// // 		hostchan <- host
// // 	}

// // 	close(hostchan)

// // 	for i := 0; i < count; i++ {
// // 		result := <-resultchan
// // 		logrus.Info(result)
// // 		if result[0] == "leader" {
// // 			leaderID = result[1]
// // 			c.ID = result[2]
// // 			return
// // 		} else if result[0] == "follower" {
// // 			if fastNodeID == "" {
// // 				fastNodeID = result[1]
// // 			}
// // 		}
// // 	}

// // 	return
// // }

// //联络人，负责联络通道中的主机，并且返回结果
// func liaison(hostchan chan string, resultchan chan []string) {

// 	for host := range hostchan {

// 		func() {
// 			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, PORT), grpc.WithInsecure())
// 			if err != nil {
// 				resultchan <- []string{"", host}
// 				return
// 			}
// 			defer conn.Close()

// 			nodeclient := rpc.NewNodeClient(conn)
// 			ctx, cancel := context.WithTimeout(context.Background(), MaxTimeout*time.Millisecond)
// 			defer cancel()
// 			result, err := nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: []byte(CMD_JOIN)})
// 			//如果发生网络错误，说明该地址下没有启动节点。
// 			if err != nil {
// 				resultchan <- []string{"", host}
// 				return
// 			}

// 			switch result.State {
// 			case 0: //follower
// 				resultchan <- []string{"follower", host}
// 			case 1: //leader
// 				resultchan <- []string{"leader", host, string(result.Data)}
// 			default:
// 			}
// 		}()

// 	}

// }
