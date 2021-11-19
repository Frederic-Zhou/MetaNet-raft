package main

import (
	"fmt"
	"metanet/node"
	"os"
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
