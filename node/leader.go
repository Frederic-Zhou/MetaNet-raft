package node

import (
	context "context"
	"fmt"
	"log"
	"time"

	"metanet/rpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//Leader is special node
type Leader = Node

//raft/rpc_call:一旦成为领导人，立即发送日志
func (l *Leader) AppendEntriesCall() {
	l.LeaderID = l.ID
	//等待迎接新节点加入

	go l.receptionNewNodes()

	// 读取出所有的节点配置地址
	for _, cfg := range l.NodesConfig {
		//初始化所有节点的 nextIndex 为自己的Log最大index+1
		cfg.NextIndex = uint64(len(l.Log))
		//向每一个节点发起链接，并 逐个推送条目
		go l.connectAndAppend(cfg)
	}
	// 检查全部Node的matchIndex ，如果大多数相同，那么,认为提交的数据成立
	for {
		//无限循环一定要加一个时间，避免占用系统处理资源。
		time.Sleep(MinTimeout / 10 * time.Millisecond)

		if l.CurrentRole != Role_Leader {
			//如果不再是Learder 退出
			return
		}

		alives := aliveNodes(l.NodesConfig)

		// logrus.Infof("all nodes count is %d", len(l.MatchIndex))
		for _, cfg := range alives {
			//规避配置还没有完全写入完毕的nil panic
			if cfg == nil {
				continue
			}

			//假设存在N
			//这个N，从节点的matchIndex中找。（这个方法是本人自己设计，而非Raft定义，Raft中没有明确定义这个N的来源）
			N := cfg.MatchIndex
			count := 0

			//满足N > commitIndex , 使得 大多数 matchIndex[i]>=N 以及 log[N].term == currentTerm
			//则令 commitIndex=N
			if N > l.CommitIndex && l.Log[N].Term == l.CurrentTerm {

				for _, cfg := range alives {
					if cfg.MatchIndex >= N {
						count++
					}
				}

				if count >= len(alives)/2 {
					l.CommitIndex = N
				}

			}
		}

	}

}

func (l *Leader) connectAndAppend(cfg *Config) {

	//*****绝对不能发给自己
	if cfg.ID == l.ID {
		return
	}

	logrus.Warnf("connect to %v self %v", cfg.ID, l.ID)

	//链接各个节点
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.ID, PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}

	defer conn.Close()

	nodeclient := rpc.NewNodeClient(conn)
	lastErr := fmt.Errorf("")

	//间隔1/3个最小随机时间，不断的给Follower发送条目或者心跳
	for {
		//*****绝对不能发给自己
		if cfg.ID == l.ID {
			break
		}
		if l.CurrentRole != Role_Leader {
			//如果不再是Learder 退出
			break
		}

		//创建一个请求参数对象，并配置其中的值

		nextIndex := cfg.NextIndex
		lastIndex := uint64(len(l.Log) - 1)

		// logrus.Infof("leaderLogLen %v,nextIndex:%v, lastIndex:%v", len(l.Log), nextIndex, lastIndex)

		//对于Follower 追加日志中尚未写入的所有条目
		entries := []*rpc.Entry{}
		//对于跟随者，最后日志条目的索引大于等于nextIndex
		if nextIndex <= lastIndex {
			entries = l.Log[nextIndex : lastIndex+1]
		}

		eArguments := &rpc.EntriesArguments{
			Term:         l.CurrentTerm,
			PrevLogIndex: nextIndex - 1,
			PrevLogTerm:  l.Log[nextIndex-1].Term,
			Entries:      entries,
			LeaderCommit: l.CommitIndex,
		}

		//创建一个超时的context，在下面进行rpc请求的时候，通过这个超时context控制请求超时
		ctx, cancel := context.WithTimeout(context.Background(), MinTimeout*time.Millisecond)
		defer cancel()
		results, err := nodeclient.AppendEntries(ctx, eArguments)
		if err != nil {
			l.nodeAliveStateChangeLog(cfg, false)

			//当新的错误与上一次错误不同时才打印，否则不打印
			if lastErr.Error() != err.Error() {
				logrus.Errorf("AppendEntries err:%v %v\n", cfg.ID, err)
			}
			lastErr = err

		} else {
			l.nodeAliveStateChangeLog(cfg, true)

			if len(entries) > 0 {
				logrus.Infof("to %v, prevlogIndex: %v, prelogTerm: %v ,commitIndex %v,lastAppliedIndex %v, entries: %v,Append: term: %v,success: %v",
					cfg.ID, eArguments.PrevLogIndex, eArguments.PrevLogTerm, l.CommitIndex, l.LastApplied, entries, results.Term, results.Success)
			}

			//如果收到的Term大于当前轮，成为Follower
			if results.Term > l.CurrentTerm {
				l.CurrentTerm = results.Term
				l.Become(Role_Follower, "心跳对象的轮比自己大")
				return
			}

			//成功后，更新对应跟随者的MatchIndex和NextIndex
			if results.Success {
				cfg.MatchIndex = lastIndex
				cfg.NextIndex = lastIndex + 1
			} else {
				//如果失败，将下一次发送日志的索引减少1，并会在次尝试发送条目
				//note: 此处做一个改善，从Follower拿到它的最新的日志index，然后从该index之后开始同步
				// cfg.NextIndex = nextIndex - 1
				cfg.NextIndex = results.FollowerLastLogIndex + 1
			}
		}

		time.Sleep(time.Millisecond * MinTimeout / 3) //1/3个最小超时时间发送一次
	}
}

func (l *Leader) receptionNewNodes() {
	//检查有没有新增的节点配置
	//如果有，发起链接和心跳
	for id := range l.NewNodeChan {
		newCfg := &Config{ID: id, NextIndex: 1, Alive: true}
		//如果添加成功与现有配置没有重复，
		if l.AddNodesConfig(newCfg) {
			go l.connectAndAppend(newCfg)
		}
	}
}

func aliveNodes(cfgs []*Config) (alives []*Config) {

	for _, cfg := range cfgs {
		if cfg.Alive {
			alives = append(alives, cfg)
		}

	}

	return

}

//根据grpc链接状态，设置节点的在线状态
//如果新状态和节点状态不同，说明状态发生改变，写入日志，同步到其他节点中去
func (l *Leader) nodeAliveStateChangeLog(cfg *Config, newState bool) {
	if cfg.Alive != newState {
		cfg.Alive = newState
		l.Log = append(l.Log, &Entry{Term: l.CurrentTerm,
			Data: []byte(fmt.Sprintf("%s%s>%t", CMD_ALIVE, cfg.ID, cfg.Alive))})
	}
}
