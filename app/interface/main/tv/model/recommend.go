package model

// ResponseRecom is the result structure from PGC API
type ResponseRecom struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Result  *ResultRecom `json:"result"`
}

// ResultRecom def.
type ResultRecom struct {
	SeasonID int      `json:"season_id"`
	From     int      `json:"from"`
	Title    string   `json:"title"`
	List     []*Recom `json:"list"`
}

// Recom def.
type Recom struct {
	Cover          string       `json:"cover"`
	FollowCount    int          `json:"follow_count"`
	IsFinish       int          `json:"is_finish"`
	IsStarted      int          `json:"is_started"`
	NewestEPCover  string       `json:"newest_ep_cover"`
	NewestEPIndex  string       `json:"newest_ep_index"`
	SeasonID       int64        `json:"season_id"`
	SeasonStatus   int          `json:"season_status"`
	SeasonType     int          `json:"season_type"`
	SeasonTypeName string       `json:"season_type_name"`
	Title          string       `json:"title"`
	TotalCount     int          `json:"total_count"`
	CornerMark     *SnVipCorner `json:"cornermark"`
}
