/*
 *
 */
package models

// Server
////////////////////////////////////////////////
// Persistent state on all servers:
// Updated on stable storage before responding to RPCs
////////////////////////////////////////////////
type Server struct {
	//Server ID
	ID string

	//Server Role, one of [Leader,Candidate,Follower]
	Role Role

	//latest term server has seen (initialized to 0 on first boot, increases monotonically)
	CurrentTerm int64

	//candidateId that received vote in current term (or null if none)
	VotedFor string

	//Log entries;
	//each entry contains command for state machine,
	//and term when entry was received by leader (first index is 1)
	//Log is all of Entries, Entries is a partment of Log
	Log []Entry

	//index of highest log entry known to be committed (initialized to 0, increases monotonically)
	CommitIndex int64

	//index of highest log entry applied to state machine (initialized to 0, increases monotonically)
	LastApplied int64

	//Volatile state on leaders:
	//(Reinitialized after election)
	//for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	NextIndex ServerIndex

	//for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
	MatchIndex ServerIndex
}

// for each server, key is server's id,value is the server's index
type ServerIndex map[string]int64

type Entry struct {
	Term int64
	Data []byte
}

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

////////////////////////////////////////////////

// AppendEntry
//Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
type AppendEntry struct {
	//leader’s term
	Term int64

	//so follower can redirect clients
	LeaderID string

	//index of log entry immediately preceding new ones
	PrevLogIndex int64

	//term of prevLogIndex entry
	PrevLogTerm int64

	//Entries is a partment of Log
	Entries []Entry

	//leader’s commitIndex
	LeaderCommit int64
}

type AppendEntryResult struct {
	//currentTerm, for leader to update itself
	Term int64

	//true if follower contained entry matching prevLogIndex and prevLogTerm
	Success bool
}

////////////////////////////////////////////////
//vote
type Vote struct {
	//candidate’s term
	Term int64

	//candidate requesting vote
	CandidateID string

	//index of candidate’s last log entry (§5.4)
	LastLogIndex int64

	//term of candidate’s last log entry (§5.4)
	LastLogTerm int64
}

type VoteResult struct {
	//currentTerm, for candidate to update itself
	Term int64

	//true means candidate received vote
	VoteGranted bool
}
