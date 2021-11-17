package node

import (
	context "context"
	"fmt"
	"log"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stefanwichmann/lanscan"
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

	//查找本地网络环境下的节点
	hosts, err := lanscan.ScanLinkLocal("tcp", PORT, 20, 5*time.Second)
	logrus.Info(hosts, err)
	if err != nil {
		return "", err
	}

	leaderID := ""
	lastNodeID := ""

	//轮训所有可连接地址
	for _, host := range hosts {
		if leaderID != "" {
			host = leaderID
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, PORT), grpc.WithInsecure())
		if err != nil {
			continue
		}

		nodeclient := rpc.NewNodeClient(conn)
		//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
		ctx, cancel := context.WithTimeout(context.Background(), MaxTimeout*time.Millisecond)
		defer cancel()
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
