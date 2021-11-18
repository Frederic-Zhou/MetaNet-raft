package node

import (
	"github.com/sirupsen/logrus"
)

func (n *Node) StateMachineHandler(data []byte) (err error) {
	logrus.Warn(string(data))
	switch {
	case string(data[:4]) == "join" && n.CurrentRole != Role_Leader: //如果是join的log条目，写入到配置，并且如果是leader通知立刻对这个节点进行连接
		id := string(data[4:])
		cfg := &Config{ID: id, NextIndex: 1}
		n.NodesConfig = append(n.NodesConfig, cfg)
	default:
		logrus.Info("Applied StateMachine:", n.Log[n.LastApplied])
	}
	return
}
