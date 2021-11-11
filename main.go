package main

import (
	"metanet/node"
)

func main() {

	n := node.NewNode("zeta")

	//开启监听
	go n.Listen()
	//加入到当前环境下的网络
	n.Join()

	n.Become(node.Role_Follower)

}
