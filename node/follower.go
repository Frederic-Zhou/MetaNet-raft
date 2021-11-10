package node

import (
	context "context"
	"metanet/rpc"
	"time"
)

type Follower = Node

var heartbeat chan byte

func (f *Follower) Timer() {
	for {
		select {
		case <-time.After(f.Timeout):
			if f.ID != f.LeaderID {
				//变成候选人
				f.RequestVoteCall()
			}
		case <-heartbeat:

		}
	}

}

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (*rpc.EntriesResults, error) {

	//收到心跳重制timer
	heartbeat <- 0

	return nil, nil
}
