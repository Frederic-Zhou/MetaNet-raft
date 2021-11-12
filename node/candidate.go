package node

import (
	context "context"
	"log"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Candidate = Node

//raft/rpc_call: Start a vote
func (c *Candidate) RequestVoteCall() bool {

	c.CurrentTerm += 1
	c.SetRandTimeOut()
	c.VotedCount = 1 //给自己投一票先
	votedTime := time.Now()

	logrus.Infof("New term is %d, timeout is %dms\n", c.CurrentTerm, c.Timeout.Milliseconds())

	for i, config := range c.NodesConfig {

		//跳过请求自己
		if config.ID == c.ID {
			continue
		}
		//初始化所有节点的 nextIndex 为自己的Log最大index+1

		c.NodesConfig[i].NextIndex = uint64(len(c.Log))
		//向每一个节点发起链接，并 逐个推送条目
		go c.connectAndVote(config)
	}

	for {

		// logrus.Infof("I'm %d, I have votedCount is %d, all node count is %d \n", c.CurrentRole, c.VotedCount, len(c.NodesConfig))
		//成为选举人
		if c.VotedCount > uint(len(c.NodesConfig)/2) {
			logrus.Infof("votedCount %d,%v", c.VotedCount, c.VotedCount > uint(len(c.NodesConfig)/2))
			return true
		}
		//超时
		if votedTime.Add(c.Timeout).Before(time.Now()) {
			//选举超时
			return false
		}
	}

}

func (c *Candidate) connectAndVote(cfg Config) {

	//链接各个节点
	conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)

	vArguments := &rpc.VoteArguments{
		Term:         c.CurrentTerm,
		CandidateID:  c.ID,
		LastLogIndex: uint64(len(c.Log) - 1),
		LastLogTerm:  c.Log[len(c.Log)-1].Term,
	}

	//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	results, err := nodeclient.RequestVote(ctx, vArguments)
	if err != nil {
		logrus.Errorf("Vote request: %v \n", err)
		return
	}

	logrus.Info(results)
	if results.VoteGranted {
		c.VotedCount += 1
	}

}
