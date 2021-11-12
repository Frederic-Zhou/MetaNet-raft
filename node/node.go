//Node is the base struct of follower, candidate, leader, client
package node

import (
	context "context"
	"fmt"
	"metanet/rpc"
	"net"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stefanwichmann/lanscan"
	"google.golang.org/grpc"
)

func NewNode() (n *Node, err error) {

	n = &Node{}

	var id uuid.UUID

	//创建ID
	for {
		id, err = uuid.NewV4()
		if err == nil {
			break
		}
	}

	n.ID = id.String()

	n.ID = strconv.Itoa(PORT) //测试

	//先填充自己的网络资料配置
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				n.Config.Address = ipnet.IP.String()
			}
		}
	}

	if n.Config.Address == "" {
		err = fmt.Errorf("no network")
		return
	}

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

func (n *Node) Work() {
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

func (f *Follower) Join() error {

	//查找本地网络环境下的节点
	hosts, err := lanscan.ScanLinkLocal("tcp4", PORT, 20, 5*time.Second)
	if err != nil {
		return err
	}

	//把自己的地址放最前面，先尝试连接自己
	hosts = append([]string{f.Config.Address}, hosts...)

	leaderAddr := ""
	for _, host := range hosts {
		if leaderAddr != "" {
			host = leaderAddr
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, PORT), grpc.WithInsecure())
		if err != nil {
			continue
		}

		nodeclient := rpc.NewNodeClient(conn)
		//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()
		result, err := nodeclient.ClientRequest(ctx, &rpc.ClientArguments{Data: []byte("join")})

		if err != nil {
			continue
		}

		if !result.Status && string(result.Data) != "" {
			leaderAddr = string(result.Data)
		}

	}

	return nil

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

	lastIndex := len(n.Log) - 1
	//如果VotedFor为空或者为CandidateID，并且候选人的日志至少和自己一样新，那么投票给他
	if n.VotedFor == "" || n.VotedFor == in.CandidateID {
		//至少一样新
		if in.LastLogIndex >= uint64(lastIndex) && in.LastLogTerm >= n.Log[lastIndex].Term {
			result.VoteGranted = true
			n.VotedFor = in.CandidateID
		}
	}

	//?????
	//note: 什么时候清理掉VotedFor的数据呢 ？

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
func (n *Node) ClientRequest(ctx context.Context, in *rpc.ClientArguments) (result *rpc.ClientResults, err error) {

	result = &rpc.ClientResults{}
	if string(in.Data) == "join" {
		result.Data = []byte{}
		result.Status = false
		if n.CurrentRole == Role_Leader {
			result.Data = []byte("ok")
			result.Status = true
		} else {
			for _, cfg := range n.NodesConfig {
				if cfg.ID == n.LeaderID {
					result.Data = []byte(cfg.Address)
					result.Status = false
					break
				}
			}
		}
	}
	return
}
