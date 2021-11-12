package main

import (
	"metanet/node"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {

	//测试=================
	//测试期间，只能创建三个节点，ID 分别是 8800、8801、8802
	node.PORT, _ = strconv.Atoi(os.Args[1])
	nodesConfig := []node.Config{
		{
			ID:        "8800",
			Address:   "127.0.0.1:8800",
			NextIndex: 1,
		},
		{
			ID:        "8801",
			Address:   "127.0.0.1:8801",
			NextIndex: 1,
		},
		{
			ID:        "8802",
			Address:   "127.0.0.1:8802",
			NextIndex: 1,
		},
	}

	//====================

	n, err := node.NewNode()
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	n.NodesConfig = nodesConfig

	//开启监听
	go n.RpcServerStart()
	//加入到当前环境下的网络
	// n.Join()

	n.Become(node.Role_Follower)
	n.Work()

}
