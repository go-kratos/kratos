package archive

const (
	// BIZVote Business Type Vote
	BIZVote = int64(2)
)

//Vote .
type Vote struct {
	VoteID    int64  `json:"vote_id"`
	VoteTitle string `json:"vote_title"`
}
