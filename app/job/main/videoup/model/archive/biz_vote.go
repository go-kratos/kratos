package archive

const (
	BIZVote = int64(2)
)

//Vote1 .
type VoteOld struct {
	VoteID    int64  `json:"vote_id"`
	VoteTitle string `json:"vote_title"`
}

//Vote .
type Vote struct {
	VoteID    int64  `json:"vote_id"`
	VoteTitle string `json:"title"`
}
