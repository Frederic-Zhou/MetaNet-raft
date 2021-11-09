package node

import (
	context "context"
	"metanet/rpc"
)

type Follower = Node

//raft/rpc_server: implemented Log and hartbeat transfer, Learder call to Followers
func (f *Follower) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (*rpc.EntriesResults, error) {
	return nil, nil
}
