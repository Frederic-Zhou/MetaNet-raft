/*
 *
 */
package models

const (
	Follower Role = iota
	Candidate
	Leader
)

// for each server, key is server's id,value is the server's index
type ServerIndex map[string]uint64

type Role int

type Entry struct {
	Term uint64
	Data []byte
}

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

	//Volatile state on leaders:
	//(Reinitialized after election)
	//for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	NextIndex ServerIndex

	//for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
	MatchIndex ServerIndex
}

////////////////////////////////////////////////

// AppendEntry
//Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
type AppendEntry struct {
	//leader’s term
	Term uint64

	//so follower can redirect clients
	LeaderID string

	//index of log entry immediately preceding new ones
	PrevLogIndex uint64

	//term of prevLogIndex entry
	PrevLogTerm uint64

	//Entries is a partment of Log
	Entries []Entry

	//leader’s commitIndex
	LeaderCommit uint64
}

type AppendEntryResult struct {
	//currentTerm, for leader to update itself
	Term uint64

	//true if follower contained entry matching prevLogIndex and prevLogTerm
	Success bool
}

////////////////////////////////////////////////
//vote
type Vote struct {
	//candidate’s term
	Term uint64

	//candidate requesting vote
	CandidateID string

	//index of candidate’s last log entry (§5.4)
	LastLogIndex uint64

	//term of candidate’s last log entry (§5.4)
	LastLogTerm uint64
}

type VoteResult struct {
	//currentTerm, for candidate to update itself
	Term uint64

	//true means candidate received vote
	VoteGranted bool
}
