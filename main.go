package main

import (
	"fmt"
	"metanet/node"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	n := node.NewNode()

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
			//加入到当前环境下的网络
			logrus.Info("尝试JOIN")
			lid, _ := n.Join()

			if lid != "" && lid != n.ID {
				logrus.Info("成功JOIN")
				return
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
		}
		result, err := n.ClientRequestCall([]byte(input))
		if err != nil {
			logrus.Error(err.Error())
		}
		logrus.Info("发送状态: %v", result.State)
	}
}
