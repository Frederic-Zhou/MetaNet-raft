package node

import (
	context "context"
	"log"
	"time"

	"metanet/rpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Leader is special node
type Leader = Node

//raft/rpc_call:
func (l *Leader) AppendEntriesCall() {

	// 读取出所有的节点配置地址
	for _, config := range l.NodesConfig {

		//初始化所有节点的 nextIndex 为自己的Log最大index+1
		l.NextIndex[config.ID] = uint64(len(l.Log))
		//向每一个节点发起链接，并 逐个推送条目
		go l.connectAndAppend(config)
	}
	// 检查全部Node的matchIndex ，如果大多数相同，那么,认为提交的数据成立
	for {
		if l.CurrentState != LeaderSTATE {
			//如果不再是Learder 退出
			break
		}

		values := []uint64{}

		// logrus.Infof("all nodes count is %d", len(l.MatchIndex))
		for _, N := range l.MatchIndex {

			count := 1

			if N > l.CommitIndex && l.Log[N].Term == l.CurrentTerm {
				for _, b := range values {
					if N == b {
						count += 1
					}
				}
				if count > len(l.MatchIndex)/2 {
					l.CommitIndex = N
				} else {
					values = append(values, N)
				}

			}
		}

	}

}

func (l *Leader) connectAndAppend(cfg Config) {

	//链接各个节点
	conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	nodeclient := rpc.NewNodeClient(conn)

	//间隔50毫秒，不断的给Follower发送条目或者心跳
	for {

		if l.CurrentState != LeaderSTATE {
			//如果不再是Learder 退出
			break
		}

		//创建一个请求参数对象，并配置其中的值
		nextIndex := l.NextIndex[cfg.ID]
		lastIndex := uint64(len(l.Log) - 1)

		entries := []*rpc.Entry{}
		if nextIndex <= lastIndex {
			entries = l.Log[nextIndex:lastIndex]
		}

		eArguments := &rpc.EntriesArguments{
			Term:         l.CurrentTerm,
			LeaderID:     l.ID,
			PrevLogIndex: nextIndex - 1,
			PrevLogTerm:  l.Log[nextIndex-1].Term,
			Entries:      entries,
			LeaderCommit: l.CommitIndex,
		}

		results, err := nodeclient.AppendEntries(ctx, eArguments)
		if err != nil {
			log.Printf("user index could not greet: %v", err)
			continue
		}

		logrus.Info(results)

		if results.Term > l.CurrentTerm {
			l.CurrentState = FollowerSTATE
		}

		//成功后，更新对应跟随者的MatchIndex和NextIndex
		if results.Success {
			l.MatchIndex[cfg.ID] = lastIndex
			l.NextIndex[cfg.ID] = lastIndex + 1
		} else {
			//如果失败，将下一次发送日志的索引减少1，并会在次尝试发送条目
			l.NextIndex[cfg.ID] = nextIndex - 1
		}

		time.Sleep(time.Millisecond * 50) //50毫秒发送一次
	}
}
