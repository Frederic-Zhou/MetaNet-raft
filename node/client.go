package node

import (
	context "context"
	"fmt"
	"log"
	"metanet/network"
	"metanet/rpc"
	"time"

	"google.golang.org/grpc"
)

type Client = Node

func (c *Client) ClientRequestCall(cmd []byte) (result *rpc.ClientResults, err error) {

	//链接到节点

	conn, err := grpc.Dial(c.LeaderID, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}

	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), MinTimeout*time.Millisecond)
	defer cancel()

	return nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: cmd})
}

func (c *Client) Join() (leaderID string, fastNodeID string) {
	//var validNetworks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

	//得到所有tcp4的网卡的网络环境的IP列表（理论列表）
	allIPs := []string{}
	for _, current := range network.LinkLocalAddresses("tcp4") {
		allIPs = append(allIPs, network.CalculateSubnetIPs(current)...)
	}

	leaderChan := make(chan string)
	nodeChan := make(chan string)

	//轮训所有可连接地址
	for _, host := range allIPs {

		go func(host string) {

			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, PORT), grpc.WithInsecure())
			if err != nil {
				nodeChan <- ""
				return
			}
			defer conn.Close()

			nodeclient := rpc.NewNodeClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), MaxTimeout*time.Millisecond)
			defer cancel()
			result, err := nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: []byte("join")})
			//如果发生网络错误，说明该地址下没有启动节点。
			if err != nil {
				nodeChan <- ""
				return
			}

			switch result.State {
			case 0: //follower
				nodeChan <- host
			case 1: //leader
				leaderChan <- host
			default:
			}

		}(host)

	}

	count := len(allIPs)
	for {
		select {
		case id := <-nodeChan:
			if fastNodeID == "" && id != "" {
				fastNodeID = id
			}

			count -= 1

			if count == 0 {
				return
			}

		case id := <-leaderChan:
			leaderID = id
			return
		}
	}

}
