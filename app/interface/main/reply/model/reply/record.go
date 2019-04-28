package reply

import (
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/xstr"
)

// time layout
const (
	RecordTimeLayout = "2006-01-02 15:04:05"
)

// Record reply record.
type Record struct {
	RpID    int64   `json:"id"`
	Oid     int64   `json:"oid"`
	Type    int32   `json:"type"`
	Floor   int32   `json:"floor"`
	Like    int32   `json:"like"`
	RCount  int32   `json:"rcount"`
	Mid     int64   `json:"mid"`
	State   int32   `json:"state"`
	Message string  `json:"message"`
	CTime   string  `json:"ctime"`
	Ats     string  `json:"ats,omitempty"`
	Members []*Info `json:"members"`
}

// FillAts fill member info of ats.
func (rc *Record) FillAts(cards map[int64]*accmdl.Card) {
	rc.Members = make([]*Info, 0, len(rc.Ats))
	ats, _ := xstr.SplitInts(rc.Ats)
	for _, at := range ats {
		if card, ok := cards[at]; ok {
			i := &Info{}
			i.FromCard(card)
			rc.Members = append(rc.Members, i)
		}
	}
	rc.Ats = ""
}
