package main

import (
	"metanet/node"
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
	}

	if lid == "" {
		logrus.Info("No Leader ... create net")
	}

	go n.RpcServerStart()
	time.Sleep(3 * time.Second)
	n.Become(node.Role_Follower)
	n.NodeWork()

}
