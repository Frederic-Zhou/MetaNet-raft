package main

import (
	"metanet/node"
)

func main() {

	n := node.NewNode("zeta")
	go n.Listen(":8800")
	n.Become(node.Role_Follower)

}
