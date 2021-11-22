package main

import (
	"bufio"
	"metanet/node"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//
func main() {
	n := node.NewNode()
	n.Become(node.Role_Follower, "启动")
	go n.ApplyStateMachine()
	go n.NodeWork()
	go n.RpcServerStart()

	simpalClient(n)
}
func simpalClient(n *node.Node) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		inputArr := strings.Split(input, " ")
		result, err := n.Command(inputArr[0], inputArr...)
		if err != nil {
			logrus.Error(err.Error())
		}
		logrus.Infof("命令返回: %v", result)
	}
}
