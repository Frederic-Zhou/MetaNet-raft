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

func (f *Follower) Join() {
	//先填充自己的网络资料配置
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
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
}

func (f *Follower) Timer() {

	for {

		select {
		case <-time.After(f.Timeout):

			logrus.Infof("I'm %d,and timeout", f.CurrentRole)

			if f.CurrentRole == Role_Follower {
				//变成候选人,发起候选投票
				f.Become(Role_Candidate)
			}
		case <-f.Heartbeat:
			logrus.Info("heartbeat,refresh timeout")
		}
	}

}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (result *rpc.EntriesResults, err error) {
	//收到心跳重制timer
	f.Heartbeat <- 0

	f.CurrentRole = Role_Follower

	result.Term = f.CurrentTerm
	result.Success = true

	//如果当前轮大于Leader的轮，返回false
	if in.Term < f.CurrentTerm {
		result.Success = false
		return
	}

	//如果没有找到匹配PrevLogIndex和PrevLogTerm的，则返回false
	if !(uint64(len(f.Log)-1) >= in.PrevLogIndex && f.Log[in.PrevLogIndex].Term == in.PrevLogTerm) {
		result.Success = false
		return
	}
	//写入到Log中
	f.Log = append(f.Log, in.GetEntries()...)

	if in.LeaderCommit > f.CommitIndex {

		if in.PrevLogIndex < in.LeaderCommit {
			f.CommitIndex = in.PrevLogIndex
		} else {
			f.CommitIndex = in.PrevLogIndex
		}

	}

	return
}

func (f *Follower) Listen() {
	// 监听本地端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", f.Port))
	if err != nil {
		fmt.Printf("listen port error: %s", err.Error())
		os.Exit(0)

	}

	// 创建gRPC服务器
	s := grpc.NewServer()
	// 注册服务
	rpc.RegisterNodeServer(s, f)
	reflection.Register(s)

	logrus.Infof("start listen %s \n", f.Port)
	err = s.Serve(lis)
	if err != nil {
		logrus.Infof("start server error: %s \n", err)
		os.Exit(0)

	}
}
