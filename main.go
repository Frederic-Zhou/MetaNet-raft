package main

import (
	"fmt"
	"metanet/node"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	n := node.NewNode()

	//加入到当前环境下的网络
	logrus.Info("尝试JOIN")

	go n.RpcServerStart()
	time.Sleep(1 * time.Second)
	n.Become(node.Role_Follower)
	go n.ApplyStateMachine()
	go n.NodeWork()

	go mustJoin(n)

	simpalClient(n)

}

func mustJoin(n *node.Node) {
	for {

		if len(n.NodesConfig) == 0 || len(n.NodesConfig) == 1 && n.NodesConfig[0].ID == n.ID {
			lid, fid := n.Join()

			if lid == "" && fid != "" {
				n.Become(node.Role_Client)
				n.LeaderID = fid
				simpalClient(n)
			}

		}

	}

}

func simpalClient(n *node.Node) {

	for {
		input := ""
		_, err := fmt.Scanln(&input)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(0)
		}
		result, err := n.ClientRequestCall([]byte(input))
		if err != nil {
			logrus.Error(err.Error())
		}
		logrus.Info(result.State)
	}
}
