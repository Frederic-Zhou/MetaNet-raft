//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"encoding/json"
	"fmt"
	"metanet/network"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
)

func NewNode() (n *Node) {

	n = &Node{}

	n.NodesConfig = []*Config{}

	//todo: 生成私钥
	n.Config.PrivateKey = []byte{}
	n.Config.PublicKey = []byte{}

	n.Log = []*rpc.Entry{
		{
			Term: 0,
			Data: []byte{},
		},
	}
	// n.Timer = time.NewTimer(0RandMillisecond())
	n.NewNodeChan = make(chan string, 20)

	return
}

func (n *Node) NodeWork() {

	for {

		switch n.CurrentRole {
		case Role_Client:

		case Role_Follower:
			n.Timer = time.NewTimer(RandMillisecond())
			<-n.Timer.C
			//Timer返回，说明超时了，身份转变为Candidate
			n.Become(Role_Candidate, "心跳超时，成为候选人")
		case Role_Candidate:
			if n.RequestVoteCall() {
				n.Become(Role_Leader, "选举获胜")
			}
		case Role_Leader:
			// 一旦成为领导人，立即发送日志
			n.LeaderID = n.ID
			n.AppendEntriesCall()
			//领导人退位的原因是收到了更高的Term
			n.Become(Role_Follower, "结束发送数据")
		}
	}
}

func (n *Node) Become(role NodeRole, reason string) {
	n.CurrentRole = role
	logrus.Infof("Change to I'm %s term is %d ,reason %v \n",
		[]string{"Client", "Follower", "Candidate", "Leader"}[n.CurrentRole],
		n.CurrentTerm, reason)
}

func (n *Node) ApplyStateMachine() {

	for {

		if n.CommitIndex > n.LastApplied {
			n.LastApplied++
			err := n.StateMachineHandler(n.Log[n.LastApplied].Data)
			if err != nil {
				n.LastApplied--
			}
		}
	}

}

//raft/rpc_server: implemented vote, after Follower change to Candidate then call to Nodes
func (n *Node) RequestVote(ctx context.Context, in *rpc.VoteArguments) (result *rpc.VoteResults, err error) {

	logrus.Warn("Receive Candidate's RequestVote...")
	//收到心跳重制timer
	n.Timer.Reset(RandMillisecond())
	logrus.Warn("Reset Timer...")

	result = &rpc.VoteResults{}
	result.Term = n.CurrentTerm
	result.VoteGranted = false

	//如果投票Term小于当前Term，返回 false
	if in.Term < n.CurrentTerm {
		return
	}

	if in.Term > n.CurrentTerm {
		n.Become(Role_Follower, "收到比自己轮大的投票")
		//如果接收到的RPC请求或响应中，任期号大于当前任期号，则当前任期号改为接收到的任期号
		n.CurrentTerm = in.Term
		n.VotedFor = ""
	}

	lastIndex := len(n.Log) - 1
	//如果VotedFor为空或者为CandidateID，并且候选人的日志至少和自己一样新，那么投票给他
	if n.VotedFor == "" || n.VotedFor == in.CandidateID {
		//至少一样新
		logrus.Warn("VotedFor ing:", n.VotedFor)
		if in.LastLogIndex >= uint64(lastIndex) && in.LastLogTerm >= n.Log[lastIndex].Term {
			result.VoteGranted = true
			n.VotedFor = in.CandidateID
			logrus.Warn("VotedFor done:", n.VotedFor)
		} else {
			logrus.Warn("VotedFor Error:", in.LastLogIndex, lastIndex, in.LastLogTerm, n.Log[lastIndex].Term)
		}
	} else {
		logrus.Warn("reject - VotedFor ing:", n.VotedFor)
	}

	return
}

//MetaNet/rpc_server: implemented node Join or Split
//*******************
//all client or new node request package to ClientRequest,include Join, rejoin, split, datarequest, ....
//all of this
func (n *Node) ClientRequest(ctx context.Context, in *rpc.ClientArguments) (result *rpc.ClientResults, err error) {
	result = &rpc.ClientResults{}

	entry := &rpc.Entry{
		Term: n.CurrentTerm,
		Data: in.Data,
	}

	//如果自己不是Leader，这种情况出现在当客户端不是节点，请求到一个不是Leader的节点时
	//不是Leader的节点请求Leader
	if n.CurrentRole != Role_Leader {
		fromid, _ := network.GetGrpcClientIP(ctx)
		ctxKV := map[string]string{
			"clientid": fromid,
		}
		return n.ClientRequestCall(in.Data, n.LeaderID, ctxKV)
	}

	//JOIN 需要立即进行ID获取，因此需要在得到请求阶段就处理，而不能等到状态机处理
	if string(in.Data) == CMD_JOIN {

		//拿到请求加入节点的地址作为ID
		fromid, selfid := network.GetGrpcClientIP(ctx)

		logrus.Infof("id %v, selfid %v, fromid %v", n.ID, selfid, fromid)
		if n.ID == "" { //如果是自己的ID是空，说明是第一次得到自己的ID，同样写入到日志中
			n.ID = selfid
			n.LeaderID = n.ID

			n.Log = append(n.Log, &rpc.Entry{
				Term: n.CurrentTerm,
				Data: []byte(fmt.Sprintf("%s%s", CMD_JOIN, selfid)),
			})
		}

		logrus.Warn("new node join:", fromid)
		//作为Leader，立即添加配置并且链接
		n.NewNodeChan <- fromid
		entry.Data = []byte(fmt.Sprintf("%s%s", CMD_JOIN, fromid))
	}

	//写入到日志中
	n.Log = append(n.Log, entry)

	result.State = 1
	return
}

func ShowNodesConfig(n *Node) {
	body, _ := json.Marshal(n.NodesConfig)
	logrus.Info(string(body))

}

func (n *Node) AddNodesConfig(newcfg *Config) (added bool) {
	add := true
	for _, cfg := range n.NodesConfig {
		if cfg.ID == newcfg.ID {
			add = false
			break
		}
	}

	if add && newcfg.ID != "" {
		n.NodesConfig = append(n.NodesConfig, newcfg)
	}

	return add && newcfg.ID != ""
}
