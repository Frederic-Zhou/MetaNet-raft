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
	"google.golang.org/grpc/reflection"
)

type Follower = Node

func (f *Follower) SetRandTimeOut() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	f.Timeout = time.Millisecond * time.Duration((MinTimeout + r.Intn(MaxTimeout-MinTimeout)))
}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (result *rpc.EntriesResults, err error) {
	logrus.Warn("============")

	//收到心跳重制timer
	f.Heartbeat.Reset(f.Timeout)

	logrus.Warn("+++++++")

	result = &rpc.EntriesResults{}
	result.Term = f.CurrentTerm
	result.Success = true

	//如果当前轮大于Leader的轮，返回false
	if in.Term < f.CurrentTerm {
		result.Success = false
		return
	}

	//如果接收到的RPC请求或响应中，任期号大于当前任期号，则当前任期号改为接收到的任期号
	f.CurrentRole = Role_Follower
	f.CurrentTerm = in.Term

	//如果没有找到匹配PrevLogIndex和PrevLogTerm的，则返回false
	if !(uint64(len(f.Log)-1) >= in.PrevLogIndex && f.Log[in.PrevLogIndex].Term == in.PrevLogTerm) {
		result.Success = false
		return
	}
	//写入到Log中
	f.Log = append(f.Log, in.GetEntries()...)

	logrus.Info(f.Log)

	//同步Leader的CommitIndex
	if in.LeaderCommit > f.CommitIndex {

		if in.PrevLogIndex < in.LeaderCommit {
			f.CommitIndex = in.PrevLogIndex
		} else {
			f.CommitIndex = in.PrevLogIndex
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
