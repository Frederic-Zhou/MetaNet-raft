package node

// for each server, key is server's id,value is the server's index
type NodeIndex map[string]uint64

//Leader is special node
type Leader struct {
	Node
	//Volatile state on leaders:
	//(Reinitialized after election)
	//for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	NextIndex NodeIndex

	//for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
	MatchIndex NodeIndex
}
