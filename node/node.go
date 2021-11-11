//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"fmt"
	"metanet/rpc"
	"net"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stefanwichmann/lanscan"
	"google.golang.org/grpc"
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
		n.Timer()
	case Role_Candidate:
		n.RequestVoteCall()
	case Role_Leader:
		n.AppendEntriesCall()

	}
}

func (f *Follower) Join() error {
	//先填充自己的网络资料配置
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				f.Config.Address = ipnet.IP.String()
			}
		}
	}

	f.Config.PrivateKey = []byte{}
	f.Config.PublicKey = []byte{}

	//查找本地网络环境下的节点
	hosts, err := lanscan.ScanLinkLocal("tcp4", PORT, 20, 5*time.Second)
	if err != nil {
		return err
	}

	for _, host := range hosts {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, PORT), grpc.WithInsecure())
		if err != nil {
			continue
		}

		nodeclient := rpc.NewNodeClient(conn)
		//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()
		result, err := nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: []byte("join")})

		if err != nil || !result.Status {
			continue
		}

	}

	return nil

}

//raft/rpc_server: implemented vote, after Follower change to Candidate then call to Nodes
func (n *Node) RequestVote(ctx context.Context, in *rpc.VoteArguments) (result *rpc.VoteResults, err error) {

	result.Term = n.CurrentTerm
	result.VoteGranted = false

	if in.Term < n.CurrentTerm {
		return
	}

	lastIndex := len(n.Log) - 1
	if n.VotedFor == "" || n.VotedFor == in.CandidateID {
		//至少一样新
		if in.LastLogIndex >= uint64(lastIndex) && in.LastLogTerm >= n.Log[lastIndex].Term {
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
