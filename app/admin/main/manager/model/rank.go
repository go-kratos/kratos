package model

// RankGroup rank permission group
type RankGroup struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Desc  string      `json:"desc"`
	Isdel int         `json:"isdel"`
	Auths []*AuthItem `json:"auths"`
}

// RankAuth rank auths
type RankAuth struct {
	ID      int64 `json:"id"`
	GroupID int64 `json:"group_id"`
	AuthID  int64 `json:"auth_id"`
	Isdel   int   `json:"isdel"`
}

// RankUser user-group-rank info.
type RankUser struct {
	ID      int64 `json:"id"`
	GroupID int64 `json:"group_id"`
	UID     int64 `json:"uid" gorm:"column:uid"`
	Rank    int   `json:"rank"`
	Isdel   int   `json:"isdel"`
}

// RankUserScores rank user scores.
type RankUserScores struct {
	UID      int64         `json:"uid"`
	Username string        `json:"username"`
	Nickname string        `json:"nickname"`
	Ranks    map[int64]int `json:"ranks"`
}
