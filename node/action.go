package node

import (
	context "context"
	"metanet/rpc"
)

func (s *Node) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (*rpc.EntriesResults, error) {
	return nil, nil
}

func (s *Node) RequestVote(ctx context.Context, in *rpc.VoteArguments) (*rpc.VoteResults, error) {
	return nil, nil
}
