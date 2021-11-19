package main

import (
	"bufio"
	"metanet/node"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	n := node.NewNode()
	n.ID = "192.168.1.100"

	go n.RpcServerStart()
	time.Sleep(1 * time.Second)
	n.Become(node.Role_Follower)
	go n.ApplyStateMachine()
	go n.NodeWork()

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
		logrus.Infof("发送状态: %v", result)
	}
}
