package node

import (
	context "context"
	"fmt"
	"math/rand"
	"metanet/rpc"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type Follower = Node

func RandMillisecond() time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Millisecond * time.Duration((MinTimeout + r.Intn(MaxTimeout-MinTimeout)))
}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (result *rpc.EntriesResults, err error) {
	logrus.Warn("Receive Leader's Append...")
	//收到心跳重制timer
	f.Timer.Reset(RandMillisecond())
	logrus.Warn("Reset Timer...")

	result = &rpc.EntriesResults{}
	result.Term = f.CurrentTerm
	result.Success = true

	//如果领导人的任期小于接收者的当前任期（接受者为Follower和Candidate）
	//note: 是否在这里判断的时候排出掉自己是Leader的情况???????
	if in.Term < f.CurrentTerm {
		result.Success = false
		return
	}

	// 所有服务器实现，如果接收到的RPC的 Term高于自己，那么更新自己的Term，并且切换为Follower
	if in.Term > f.CurrentTerm {
		f.Become(Role_Follower)
		//如果接收到的RPC请求或响应中，任期号大于当前任期号，则当前任期号改为接收到的任期号
		f.CurrentTerm = in.Term
		f.VotedFor = ""
	}

	//如果没有找到匹配PrevLogIndex和PrevLogTerm的，则返回false
	//note: uint64(len(f.Log)-1) >= in.PrevLogIndex 这个条件主要是规避数组越界问题。Raft的要求 && 后的第二个判断即可满足
	if !(uint64(len(f.Log)-1) >= in.PrevLogIndex && f.Log[in.PrevLogIndex].Term == in.PrevLogTerm) {
		logrus.Errorf("第一个条件为%v,第二个条件为%v", uint64(len(f.Log)-1) >= in.PrevLogIndex, f.Log[in.PrevLogIndex].Term == in.PrevLogTerm)
		result.Success = false
		return
	}

	//如果一个已经存在的条目和新条目冲突（索引相同，任期不同）删除该条目即之后的条目
	if uint64(len(f.Log)-1) >= in.PrevLogIndex+1 && f.Log[in.PrevLogIndex].Term != in.PrevLogTerm {
		f.Log = f.Log[:in.PrevLogIndex+1]
	}

	//追加日志中尚未存在的任何条目
	f.Log = append(f.Log, in.GetEntries()...)

	// logrus.Info(f.Log)

	//同步Leader的CommitIndex
	//note: 有可能出现本节点是较慢的节点，Leader已经提交了较高的Index，但是本节点未必获得了领导的最高提交Index的日志
	//因此，如果本节点的上一条日志高于leader Commited的日志，那么就用leader Commited的最高日志
	//如果，Leader提交的最高日志比本节点上一次接受到的高，说明本节点还没有接收到leader提交的最高日志以及中间的日志，那么本节点的提交日志实用PrevLogIndex
	if in.LeaderCommit > f.CommitIndex {

		if in.PrevLogIndex < in.LeaderCommit {
			f.CommitIndex = in.PrevLogIndex
		} else {
			f.CommitIndex = in.LeaderCommit
		}

	}

	// f.LeaderID = in.LeaderID
	// 动态节点，将发心跳过来的IP 作为 LeaderID
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			f.LeaderID = tcpAddr.IP.String()
		} else {
			f.LeaderID = pr.Addr.String()
		}
	}

	return
}

func (f *Follower) RpcServerStart() {
	// 监听本地端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		fmt.Printf("listen port error: %s", err.Error())
		os.Exit(0)

	}

	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	rpc.RegisterNodeServer(s, f)
	reflection.Register(s)

	logrus.Infof("start listen \n")
	err = s.Serve(lis)
	if err != nil {
		logrus.Infof("start server error: %s \n", err)
		os.Exit(0)

	}
}
