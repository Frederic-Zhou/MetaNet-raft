package node

import (
	context "context"
	"metanet/rpc"
)

// Server change name to Node
////////////////////////////////////////////////
// Persistent state on all servers:
// Updated on stable storage before responding to RPCs
////////////////////////////////////////////////
type Node struct {
	//latest term server has seen (initialized to 0 on first boot, increases monotonically)
	CurrentTerm uint64

	//candidateId that received vote in current term (or null if none)
	VotedFor string

	//Log entries;
	//each entry contains command for state machine,
	//and term when entry was received by leader (first index is 1)
	//Log is all of Entries, Entries is a partment of Log
	Log []Entry

	//index of highest log entry known to be committed (initialized to 0, increases monotonically)
	CommitIndex uint64

	//index of highest log entry applied to state machine (initialized to 0, increases monotonically)
	LastApplied uint64

	//self extention informations
	Extention

	//other known node extentions
	NodesExtention []Extention
}

type Extention struct {
	//Node ID
	ID string
	//use to display name
	Name string
	//Node network address
	Address    string
	PrivateKey string
	PublicKey  string
}

type Entry struct {
	Term uint64
	Data []byte
}

//raft: implemented Log and hartbeat transfer, Learder call to Followers
func (n *Node) AppendEntries(ctx context.Context, in *rpc.EntriesArguments) (*rpc.EntriesResults, error) {
	return nil, nil
}

//raft: implemented vote, after Follower change to Candidate then call to Nodes
func (n *Node) RequestVote(ctx context.Context, in *rpc.VoteArguments) (*rpc.VoteResults, error) {
	return nil, nil
}

//MetaNet: implemented node Join or Split
//the node need to input:
//1. network address,include IP、 bluethooth
//2. if the node first time Join to metanet, Node must input passphrase, and send publicKey to metanet
//   if the node Rejoin to metanet, check the signtrue instead of input passphrase.
//3. there are three action in Join function:
//   a. join
//	 b. rejoin
//   c. split
func (n *Node) Join() {}
