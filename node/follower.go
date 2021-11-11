package node

import (
	context "context"
	"math/rand"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
)

type Follower = Node

var heartbeat chan byte

func (f *Follower) SetRandTimeOut() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	f.Timeout = time.Millisecond * time.Duration((MinTimeout + r.Intn(MaxTimeout-MinTimeout)))
}

func (f *Follower) Timer() {

	for {

		select {
		case <-time.After(f.Timeout):
			logrus.Infof("I'm %d,and timeout", f.CurrentState)
			if f.CurrentState == FollowerSTATE {
				//变成候选人,发起候选投票
				f.ChangeTo(CandidateSTATE)
			}
		case <-heartbeat:
			logrus.Info("heartbeat,refresh timeout")
		}
	}

}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (result *rpc.EntriesResults, err error) {
	//收到心跳重制timer
	heartbeat <- 0

	f.CurrentState = FollowerSTATE

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
