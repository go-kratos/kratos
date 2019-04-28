package thirdp

import (
	"fmt"

	"go-common/library/time"
)

// MangoRecom is mango recom table structure
type MangoRecom struct {
	ID        int64     `json:"id"`
	RID       int64     `json:"rid"`
	Rtype     int       `json:"rtype"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Category  int       `json:"category"`
	Playcount int64     `json:"playcount"`
	JID       int64     `json:"jid"`
	Content   string    `json:"content"`
	Staff     string    `json:"staff"`
	Rorder    int       `json:"rorder"`
	Mtime     time.Time `json:"-"`
}

// MangoParams is the output structure for mango recom api
type MangoParams struct {
	JumpParam string `json:"jump_param"` // combine
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Playcount int64  `json:"play_count"`
	Category  string `json:"category"` // transform to CN
	Desc      string `json:"desc"`
	Staff     string `json:"staff"`
	Role      string `json:"role"`      // from DB
	PlayTime  string `json:"play_time"` // from DB
}

// MangoOrder struct
type MangoOrder struct {
	RIDs []int64
}

const (
	_fixStr   = "&progress=0&from=mango&resource=rec"
	_pgcJump  = "yst://com.xiaodianshi.tv.yst?type=3&isBangumi=1&seasonId=%d&epId=%d" + _fixStr
	_ugcJump  = "yst://com.xiaodianshi.tv.yst?type=3&isBangumi=0&avId=%d&cId=%d" + _fixStr
	_rtypePGC = 1
	_rtypeUGC = 2
)

// ToParam transforms an MangoRecom from DB to MangoParam for mango OS
func (m *MangoRecom) ToParam() *MangoParams {
	param := &MangoParams{
		Title:     m.Title,
		Cover:     m.Cover,
		Playcount: m.Playcount,
		Desc:      m.Content,
		Staff:     m.Staff,
	}
	if m.Rtype == _rtypePGC {
		param.JumpParam = fmt.Sprintf(_pgcJump, m.RID, m.JID)
	}
	if m.Rtype == _rtypeUGC {
		param.JumpParam = fmt.Sprintf(_ugcJump, m.RID, m.JID)
	}
	return param
}
