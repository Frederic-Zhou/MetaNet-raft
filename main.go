package main

import (
	"metanet/node"
)

func main() {

	n := node.NewNode("zeta")

	n.Become(node.Role_Follower)

}
