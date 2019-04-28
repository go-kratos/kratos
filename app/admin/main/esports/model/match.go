package model

// Match .
type Match struct {
	ID       int64  `json:"id" form:"id"`
	Title    string `json:"title" form:"title"`
	SubTitle string `json:"sub_title" form:"sub_title"`
	CYear    int    `json:"c_year" form:"c_year"`
	Sponsor  string `json:"sponsor" form:"sponsor"`
	Logo     string `json:"logo" form:"logo" validate:"required"`
	Dic      string `json:"dic" form:"dic"`
	Status   int    `json:"status" form:"status"`
	Rank     int    `json:"rank" form:"rank" validate:"min=0,max=99"`
}

// MatchInfo .
type MatchInfo struct {
	*Match
	Games []*Game `json:"games"`
}

// TableName .
func (m Match) TableName() string {
	return "es_matchs"
}
