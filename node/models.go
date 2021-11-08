package node

//**Node is Server nickname**

// for each server, key is server's id,value is the server's index
type NodeIndex map[string]uint64

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Entry struct {
	Term uint64
	Data []byte
}

// Server change name to Node
////////////////////////////////////////////////
// Persistent state on all servers:
// Updated on stable storage before responding to RPCs
////////////////////////////////////////////////
type Node struct {
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
	NextIndex NodeIndex

	//for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
	MatchIndex NodeIndex

	Extentions
}

type Extentions struct {
	Nickname   string
	Address    string
	PrivateKey string
	PublicKey  string
}
