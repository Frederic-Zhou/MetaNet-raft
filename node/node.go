//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"metanet/rpc"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
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
	n.Port = 8800

	n.SetRandTimeOut()
	return
}

func (n *Node) Become(role NodeRole) {
	n.CurrentRole = role

	logrus.Infof("Change to %d \n", role)

	logrus.Infof("Now I'm %d term is %d, timeout is %dms \n",
		n.CurrentRole,
		n.CurrentTerm,
		n.Timeout.Milliseconds())

	switch role {
	case Role_Client:

	case Role_Follower:
		//开启监听
		go n.Listen()
		//加入到当前环境下的网络
		n.Join()
		//记时
		n.Timer()
	case Role_Candidate:
		n.RequestVoteCall()
	case Role_Leader:
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
