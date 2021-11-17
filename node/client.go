package node

import (
	context "context"
	"fmt"
	"log"
	"metanet/network"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
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

func (c *Client) Join() (string, error) {
	//var validNetworks = []string{"tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "ip", "ip4", "ip6", "unix", "unixgram", "unixpacket"}

	//得到所有tcp4的网卡的网络环境的IP列表（理论列表）
	allIPs := []string{}
	for _, current := range network.LinkLocalAddresses("tcp4") {
		allIPs = append(allIPs, network.CalculateSubnetIPs(current)...)
	}

	leaderID := ""
	lastNodeID := ""
	//轮训所有可连接地址
	for _, host := range allIPs {
		if leaderID != "" {
			host = leaderID
		}
		logrus.Info("try:", host)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, PORT), grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			continue
		}

		nodeclient := rpc.NewNodeClient(conn)

		result, err := nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: []byte("join")})

		logrus.Info("Join ", host, err, result)
		//如果发生网络错误，说明该地址下没有启动节点。
		if err != nil {
			//如果leaderID 是存在的，说明此节点无法链接Leader，直接跳出
			if leaderID != "" {
				break
			}
			continue
		}

		//找到的是follower节点，返回leaderID ,下一次连接尝试链接Leader
		if result.State == 0 {
			leaderID = string(result.Data)
			lastNodeID = host
			continue
		}

		//找到了Leader，并已经成功加入, Leader返回本节点ID（IP 地址）
		if result.State == 1 {
			return string(result.Data), nil
		}
	}

	return lastNodeID, fmt.Errorf("join fail")

}
