package node

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

func (n *Node) StateMachineHandler(data []byte) (err error) {
	switch {
	//如果是join的log条目，写入到配置，并且如果是leader 已经添加，无需再配置通知立刻对这个节点进行连接
	case bytes.HasPrefix(data, []byte(CMD_JOIN)) && n.CurrentRole != Role_Leader:
		id := string(data[4:])
		cfg := &Config{ID: id, NextIndex: 1}
		n.AddNodesConfig(cfg)
	default:
		logrus.Info("Applied StateMachine:", n.Log[n.LastApplied])
	}

	// jsondata, _ := json.Marshal(n.Log)
	// logrus.Info(string(jsondata))

	// jsondata, _ = json.Marshal(n.NodesConfig)
	// logrus.Info(string(jsondata))

	return
}
