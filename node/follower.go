package node

import (
	context "context"
	"fmt"
	"math/rand"
	"metanet/network"
	"metanet/rpc"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Follower = Node

func RandMillisecond() time.Duration {
	// logrus.Error("set了一次")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Millisecond * time.Duration((MinTimeout + r.Intn(MaxTimeout-MinTimeout)))
}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (result *rpc.EntriesResults, err error) {
	//收到心跳重制timer
	f.Timer.Reset(RandMillisecond())

	// rid, sid := network.GetGrpcClientIP(ctx)
	// logrus.Warnf("---%v,%v", rid, sid)

	result = &rpc.EntriesResults{}
	result.Term = f.CurrentTerm
	result.Success = true

	// 所有服务器实现，如果接收到的RPC的 Term高于自己，那么更新自己的Term，并且切换为Follower
	if in.Term > f.CurrentTerm {

		//如果接收到的RPC请求或响应中，任期号大于当前任期号，则当前任期号改为接收到的任期号
		f.CurrentTerm = in.Term
		f.VotedFor = ""
		f.Become(Role_Follower, "收到一个比自己轮大的心跳")
	}

	//此判断语句以下的实现都要排除自己是Leader的情况
	//note:以下判断是错误的，无论如何都不能接收，尽管直接返回，但是不管result.Success是true还是false，都会影响集群结果
	//所以注释掉，在Leader发送数据的地方处理，绝对不能给自己发
	// if f.CurrentRole == Role_Leader {
	// 	return
	// }
	// logrus.Warning("ROLE:", f.CurrentRole)

	//如果领导人的任期小于接收者的当前任期（接受者为Follower和Candidate）
	//note: 是否在这里判断的时候排出掉自己是Leader的情况???? :针对这个问题，在Leader发送时做了判断，一定不能给自己发。
	if in.Term < f.CurrentTerm {
		result.Success = false
		return
	}

	lastIndex := uint64(len(f.Log) - 1)
	//如果没有找到匹配PrevLogIndex和PrevLogTerm的，则返回false
	//note: uint64(len(f.Log)-1) >= in.PrevLogIndex 这个条件主要是规避数组越界问题。Raft的要求 && 后的第二个判断即可满足
	if !(lastIndex >= in.PrevLogIndex && f.Log[in.PrevLogIndex].Term == in.PrevLogTerm) {
		logrus.Errorf("in.PrevLogIndex %v", in.PrevLogIndex)
		// logrus.Errorf("第一个条件为%v,第二个条件为%v",
		// 	uint64(len(f.Log)-1) >= in.PrevLogIndex,
		// 	f.Log[in.PrevLogIndex].Term == in.PrevLogTerm)
		result.Success = false
		result.FollowerLastLogIndex = lastIndex
		return
	}

	//如果一个已经存在的条目和新条目（译者注：即刚刚接收到的日志条目）发生了冲突（因为索引相同，任期不同），
	//那么就删除这个已经存在的条目以及它之后的所有条目 （5.3 节）
	//note:这里的判断似乎和论文上有差异，还需要修正????
	if lastIndex >= in.PrevLogIndex+1 && f.Log[in.PrevLogIndex+1].Term != in.Entries[0].Term {
		f.Log = f.Log[:in.PrevLogIndex+1]
	}

	//追加日志中尚未存在的任何条目
	f.Log = append(f.Log, in.Entries...)

	if len(in.Entries) > 0 {
		logrus.Info("new log is:", in.Entries)
	}
	// 动态节点，将发心跳过来的IP 作为 LeaderID, 请求目标IP就是自己的ID
	f.LeaderID, f.ID = network.GetGrpcClientIP(ctx)

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

	if len(in.Entries) > 0 {
		logrus.Warnf("leaderCommit %v, in.PrevlogIndex %v,commitindex %v, lastapplied %v", in.LeaderCommit, in.PrevLogIndex, f.CommitIndex, f.LastApplied)
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
