package model

import (
	"go-common/library/time"
)

const (
	_TypeUGC = 2
	_TypePGC = 1
)

// Homepage is the home page struct
type Homepage struct {
	Recom  []*Card            `json:"recom"`
	Latest []*Card            `json:"latest"`
	Lists  map[string][]*Card `json:"lists"`
	Follow []*Follow          `json:"follow,omitempty"`
}

// Card is the unit to display
type Card struct {
	SeasonID   int          `json:"season_id"`
	Title      string       `json:"title"`
	Cover      string       `json:"cover"`
	Type       int          `json:"type"` // 1=pgc, 2=ugc
	NewEP      *NewEP       `json:"new_ep"`
	CornerMark *SnVipCorner `json:"cornermark"`
}

// IsUGC returns whether the card is ugc card
func (c Card) IsUGC() bool {
	return c.Type == _TypeUGC
}

// BePGC def.
func (c *Card) BePGC() {
	c.Type = _TypePGC
}

// NewEP is the latest EP of a season
type NewEP struct {
	ID        int64  `json:"id"`
	Index     string `json:"index"`
	IndexShow string `json:"index_show"`
	Cover     string `json:"cover"`
}

// Rank represents the table TV_RANK
type Rank struct {
	ID        int64
	Rank      int
	Title     string
	Type      int8
	CID       int64
	ContID    int64
	Category  int8
	Position  int32
	IsDeleted int8
	Ctime     time.Time
	Mtime     time.Time
}

// SimpleRank picks the necessary fields from tv_rank
type SimpleRank struct {
	ContID   int64
	ContType int
}

// RespModInterv is the response struct for mod intervention
type RespModInterv struct {
	Ranks []*SimpleRank
	AIDs  []int64
	SIDs  []int64
}

// IsUGC returns whether the card is ugc card
func (c SimpleRank) IsUGC() bool {
	return c.ContType == _TypeUGC
}

// ReqZone is the request struct of zone page
type ReqZone struct {
	SType       int
	IntervType  int
	LengthLimit int
	IntervM     int
	PGCListM    map[int][]*Card
}

//RespAI is the response of AI ugc rank data
type RespAI struct {
	Note       string    `json:"note"`
	SourceData string    `json:"source_data"`
	Code       int       `json:"code"`
	Num        int       `json:"num"`
	List       []*AIData `json:"list"`
}

// AIData is the ai card structure
type AIData struct {
	AID         int `json:"aid"`
	MID         int `json:"mid"`
	Pts         int `json:"pts"`
	Play        int `json:"play"`
	Coints      int `json:"coins"`
	VideoReview int `json:"video_review"`
}

//ToCard transforms an ArcCMS to Card
func (a ArcCMS) ToCard() *Card {
	return &Card{
		SeasonID: int(a.AID),
		Title:    a.Title,
		Cover:    a.Cover,
		Type:     _TypeUGC,
		NewEP:    &NewEP{Cover: a.Cover},
	}
}

//ToIdxSn transforms an ArcCMS to IdxSeason
func (a ArcCMS) ToIdxSn() *IdxSeason {
	return &IdxSeason{
		SeasonID: a.AID,
		Title:    a.Title,
		Cover:    a.Cover,
		Upinfo:   "",
	}
}

// ReqZoneInterv is the request structure for zone intervention
type ReqZoneInterv struct {
	RankType int
	Category int
	Limit    int
}
