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
	if id, err := n.Join(); err != nil && id != "" {
		//如果加入失败,保存可发送请求的ID
		//不启动节点功能
		n.Become(node.Role_Client)
		n.LeaderID = id

	} else {
		//开启监听，并启动节点功能
		go n.RpcServerStart()
		time.Sleep(3 * time.Second)
		n.Become(node.Role_Follower)
		n.NodeWork()
	}

}
