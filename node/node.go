//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"metanet/rpc"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

func NewNode() (n *Node, err error) {

	n = &Node{}
	n.Config.PrivateKey = []byte{}
	n.Config.PublicKey = []byte{}

	n.Log = []*rpc.Entry{
		{
			Term: 0,
			Data: []byte{},
		},
	}

	n.Timer = time.NewTimer(RandMillisecond())

	return
}

func (n *Node) NodeWork() {
	for {
		switch n.CurrentRole {
		case Role_Client:

		case Role_Follower:
			<-n.Timer.C
			//Timer返回，说明超时了，身份转变为Candidate
			n.Become(Role_Candidate)
		case Role_Candidate:
			if n.RequestVoteCall() {
				n.Become(Role_Leader)
			}
		case Role_Leader:
			// 一旦成为领导人，立即发送日志
			n.AppendEntriesCall()
			//领导人退位的原因是收到了更高的Term
			n.Become(Role_Follower)
		}
	}
}

func (n *Node) Become(role NodeRole) {
	n.CurrentRole = role
	logrus.Infof("Now I'm %d term is %d \n",
		n.CurrentRole,
		n.CurrentTerm)
}

func (n *Node) ApplyStateMachine() {
	//todo: 应用到状态机
	if n.CommitIndex > n.LastApplied {
		n.LastApplied++
		logrus.Debug("...", n.Log[n.LastApplied])
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
		n.Become(Role_Follower)
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

	//如果收到接入请求
	if string(in.Data) == "join" {

		//如果自己是Leader
		if n.CurrentRole == Role_Leader {
			id := ""
			if pr, ok := peer.FromContext(ctx); ok {
				if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
					id = tcpAddr.IP.String()
				} else {
					id = pr.Addr.String()
				}
			}

			result.State = 1
			result.Data = []byte(n.ID)

			cfg := Config{NextIndex: 1, ID: id}
			n.NodesConfig = append(n.NodesConfig, cfg)

			logrus.Warn("new node join", n.NodesConfig)
			go n.connectAndAppend(&cfg)
		} else {
			result.State = 0
			result.Data = []byte(n.LeaderID)
		}

		return
	}

	//如果自己不是Leader，调用自己的请求，转发给Leader，这种情况出现在当客户端不是节点，请求到一个不是Leader的节点时
	//不是Leader的节点用自己的请求函数去请求Leader
	if n.CurrentRole != Role_Leader {
		return n.ClientRequestCall(in.Data)
	}

	entry := &rpc.Entry{
		Term: n.CurrentTerm,
		Data: in.Data,
	}

	//写入到日志中
	n.Log = append(n.Log, entry)

	result.State = 1
	return
}
