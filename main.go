package main

import (
	"fmt"
	"metanet/node"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	n, err := node.NewNode()
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	//加入到当前环境下的网络
	logrus.Info("尝试JOIN")
	lid, fid := n.Join()

	if lid == "" && fid != "" {
		n.Become(node.Role_Client)
		n.LeaderID = fid
		simpalClient(n)
	}

	if lid == "" {
		logrus.Info("No Leader ... create net")
	}

	go n.RpcServerStart()
	time.Sleep(3 * time.Second)
	n.Become(node.Role_Follower)
	go n.ApplyStateMachine()
	go n.NodeWork()

	simpalClient(n)

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
