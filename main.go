package main

import (
	"metanet/node"
)

func main() {

	n := node.NewNode("zeta")
	go n.Listen(":8800")
	n.ChangeTo(node.FollowerSTATE)

}
