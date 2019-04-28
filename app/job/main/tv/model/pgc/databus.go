package pgc

// MediaEP is the new structure of ep in Databus Msg
type MediaEP struct {
	ID        int64  `json:"id"`
	EPID      int    `json:"epid"`
	SeasonID  int    `json:"season_id"`
	State     int    `json:"state"`
	Valid     int    `json:"valid"`
	IsDeleted int    `json:"is_deleted"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Cover     string `json:"cover"`
	Mark      int    `json:"mark"`
	CID       int64  `json:"cid"`
	PayStatus int    `json:"pay_status"`
}

// MediaSn is the new structure of season in Databus Msg
type MediaSn struct {
	ID          int64  `json:"id"`
	IsDeleted   int8   `json:"is_deleted"`
	Valid       int    `json:"valid"`
	Check       int8   `json:"check"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Desc        string `json:"desc"`
	UpInfo      string `json:"upinfo"`
	Ctime       string `json:"ctime"`
	Category    int    `json:"category"`
	Area        string `json:"area"`
	Playtime    string `json:"play_time"`
	Role        string `json:"role"`
	Staff       string `json:"staff"`
	TotalNum    int    `json:"total_num"`
	Style       string `json:"style"`
	Producer    string `json:"producer"`
	Version     string `json:"version"`
	AliasSearch string `json:"alias_search"`
	Brief       string `json:"brief"`
	Status      int    `json:"status"`
}

// DatabusRes is the result of databus message
type DatabusRes struct {
	Action string `json:"action"`
	Table  string `json:"table"`
}

// DatabusEP is the struct of message for the modification of tv_content
type DatabusEP struct {
	New *MediaEP `json:"new"`
	Old *MediaEP `json:"old"`
}

// DatabusSeason is the struct of message for the modification of tv_ep_season
type DatabusSeason struct {
	Old *MediaSn `json:"old"`
	New *MediaSn `json:"new"`
}

// ToSimple returns SimpleSeason struct
func (m *MediaSn) ToSimple() *SimpleSeason {
	return &SimpleSeason{
		ID:        m.ID,
		IsDeleted: m.IsDeleted,
		Valid:     m.Valid,
		Check:     m.Check,
	}
}

// ToSimple returns SimpleEP struct
func (ep *MediaEP) ToSimple() *SimpleEP {
	return &SimpleEP{
		ID:        ep.ID,
		IsDeleted: ep.IsDeleted,
		Valid:     ep.Valid,
		State:     ep.State,
		SeasonID:  ep.SeasonID,
		EPID:      ep.EPID,
		NoMark:    ep.Mark,
	}
}

// ToCMS returns EpCMS
func (ep *MediaEP) ToCMS() *EpCMS {
	return &EpCMS{
		EPID:      int(ep.EPID),
		Cover:     ep.Cover,
		Title:     ep.Title,
		Subtitle:  ep.Subtitle,
		PayStatus: ep.PayStatus,
	}
}
