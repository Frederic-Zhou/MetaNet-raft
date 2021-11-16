//Node is the base struct of follower, candidate, leader, client
package node

import (
	"metanet/rpc"
	"sync"
	"time"
)

type NodeRole uint8

const (
	Role_Client NodeRole = iota
	Role_Follower
	Role_Candidate
	Role_Leader

	MinTimeout = 150 * 10
	MaxTimeout = 300 * 10
)

var PORT = 8800

type Config struct {
	//Node ID
	ID string

	PrivateKey []byte
	PublicKey  []byte

	CurrentRole NodeRole

	//下一次应该发送的索引位置
	NextIndex uint64
	//已复制的最大Log索引
	MatchIndex uint64
}

// type Entry struct {
// 	// Index uint64
// 	Term uint64
// 	Data []byte
// }
//直接使用RPC中定义的Entry，减少类型转换
type Entry = rpc.Entry

// for each server, key is server's id,value is the server's index
type NodeIndex = sync.Map

////////////////////////////////////////////////
// Persistent state on all servers:
// Updated on stable storage before responding to RPCs
////////////////////////////////////////////////
type Node struct {
	//latest term server has seen (initialized to 0 on first boot, increases monotonically)
	//当前轮的编号，Leader用此编号发送条目，以便于与Follower沟通数据一致性问题
	//Candidate 会自增轮号
	CurrentTerm uint64

	//candidateId that received vote in current term (or null if none)
	VotedFor string

	//current term LeaderID
	LeaderID string

	//
	VotedCount uint

	//Log entries;
	//each entry contains command for state machine,
	//and term when entry was received by leader (first index is 1)
	//Log is all of Entries, Entries is a partment of Log
	Log []*Entry

	//index of highest log entry known to be committed (initialized to 0, increases monotonically)
	CommitIndex uint64

	//index of highest log entry applied to state machine (initialized to 0, increases monotonically)
	LastApplied uint64

	//Volatile state on leaders:
	//(Reinitialized after election)
	//for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	// NextIndex NodeIndex //包含到NodesConfig中

	//for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
	// MatchIndex NodeIndex  //包含到NodesConfig中

	//self Config informations
	Config

	//other known nodes configs
	NodesConfig []Config

	Timer *time.Timer
}
