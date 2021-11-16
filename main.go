package main

import (
	"metanet/node"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	//测试=================
	//测试期间，只能创建三个节点，ID 分别是 8800、8801、8802
	// node.PORT, _ = strconv.Atoi(os.Args[1])
	// nodesConfig := []node.Config{
	// 	{
	// 		ID:        ":8800",
	// 		NextIndex: 1,
	// 	},
	// 	{
	// 		ID:        ":8801",
	// 		NextIndex: 1,
	// 	},
	// 	{
	// 		ID:        ":8802",
	// 		NextIndex: 1,
	// 	},
	// }

	//====================

	n, err := node.NewNode()
	if err != nil {
		logrus.Errorln(err.Error())
		return
	}

	// n.NodesConfig = nodesConfig

	//加入到当前环境下的网络
	logrus.Info("尝试JOIN")
	if id, err := n.Join(); err != nil && id != "" {
		//如果加入失败,保存可发送请求的ID
		n.Become(node.Role_Client)
		n.LeaderID = id

	} else {
		//开启监听
		go n.RpcServerStart()
		time.Sleep(3 * time.Second)
		n.Become(node.Role_Follower)
		n.RaftWork()
	}

}
