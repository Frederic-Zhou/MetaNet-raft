package node

import (
	context "context"
	"fmt"
	"log"
	"metanet/rpc"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Candidate = Node

//raft/rpc_call: Start a vote
func (c *Candidate) RequestVoteCall() bool {

	c.CurrentTerm += 1                         //当前任期自增
	c.VotedCount = 1                           //先给自己记录一票
	c.Timer = time.NewTimer(RandMillisecond()) //启动一个定时器

	logrus.Infof("New term is %d \n", c.CurrentTerm)

	alives := aliveNodes(c.NodesConfig)

	for _, cfg := range alives {
		//排除自己
		if cfg.ID == c.ID {
			continue
		}

		//初始化所有节点的 nextIndex 为自己的Log最大index+1
		cfg.NextIndex = uint64(len(c.Log))
		//向每一个节点发起链接，并 逐个推送条目
		go c.connectAndVote(cfg)
	}

	for {
		select {
		case <-c.Timer.C:
			return false
		default:
			// case <-time.After(10 * time.Millisecond):
			//如果收到大多数服务器的选票，成为领导人
			// logrus.Infof("votedCount %d/%d at %d \n", c.VotedCount, len(c.NodesConfig), c.CurrentTerm)
			if c.VotedCount > uint(len(alives)/2) {
				c.Timer.Stop()
				return true
			}

			time.Sleep(MinTimeout / 10 * time.Millisecond)

		}

	}

}

func (c *Candidate) connectAndVote(cfg *Config) {

	//链接各个节点
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.ID, PORT), grpc.WithInsecure())

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*MinTimeout)
	defer cancel()
	results, err := nodeclient.RequestVote(ctx, vArguments)
	if err != nil {
		logrus.Errorf("Vote request: %v \n", err)
		return
	}

	if results.Term > c.CurrentTerm {
		c.CurrentTerm = results.Term
		c.Become(Role_Follower, "投票人的轮数比自己还大")
		return
	}

	logrus.Info("vote return:", results)
	if results.VoteGranted {
		c.VotedCount++
	}

}
