//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"fmt"
	"metanet/rpc"
	"net"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewNode(name string) (n *Node) {

	n = &Node{}

	var id uuid.UUID
	var err error
	for {
		id, err = uuid.NewV4()
		if err == nil {
			break
		}
	}

	n.ID = id.String()

	n.SetRandTimeOut()
	return
}

func (n *Node) ChangeTo(state NodeState) {
	n.CurrentState = state

	logrus.Infof("Change to %d \n", state)

	logrus.Infof("Now I'm %d term is %d, timeout is %dms \n",
		n.CurrentState,
		n.CurrentTerm,
		n.Timeout.Milliseconds())
	switch state {
	case ClientSTATE:

	case FollowerSTATE:
		n.Timer()
	case CandidateSTATE:
		n.RequestVoteCall()
	case LeaderSTATE:
		n.AppendEntriesCall()

	}
}

//raft/rpc_server: implemented vote, after Follower change to Candidate then call to Nodes
func (n *Node) RequestVote(ctx context.Context, in *rpc.VoteArguments) (result *rpc.VoteResults, err error) {

	result.Term = n.CurrentTerm
	result.VoteGranted = false

	if in.Term < n.CurrentTerm {
		return
	}

	if n.VotedFor == "" || n.VotedFor == in.CandidateID {
		//至少一样新
		if in.LastLogIndex >= uint64(len(n.Log)-1) {
			result.VoteGranted = true
		}
	}

	return
}

//MetaNet/rpc_server: implemented node Join or Split
//the node need to input:
//1. network address,include IP、 bluethooth
//2. if the node first time Join to metanet, Node must input passphrase, and send publicKey to metanet
//   if the node Rejoin to metanet, check the signtrue instead of input passphrase.
//3. there are three action in Join function:
//   a. join
//	 b. rejoin
//   c. split

//*******************
//all client or new node request package to ClientRequest,include Join, rejoin, split, datarequest, ....
//all of this
func (n *Node) ClientRequest(ctx context.Context, in *rpc.ClientArguments) (*rpc.ClientResults, error) {
	return nil, nil
}

func (n *Node) Listen(port string) {

	// 监听本地端口
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}

	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	rpc.RegisterNodeServer(s, n)
	reflection.Register(s)

	logrus.Info("开始监听")
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}
}
