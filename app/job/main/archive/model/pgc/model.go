package pgc

const (
	TypeForBangumi     = 1
	TypeForMovie       = 2
	TypeForDocumentary = 3
)

type MovieEpisode struct {
	EpID     int64 `json:"id"`
	SeasonID int64 `json:"movieSeasonId"`
	AID      int64 `json:"aid"`
	CID      int64 `json:"cid"`
	Status   int8  `json:"status"`
}

type BangumiEpisode struct {
	EpID     int64 `json:"episodeId"`
	AID      int64 `json:"avId"`
	SeasonID int64 `json:"seasonId"`
	CID      int64 `json:"danmaku"`
	Status   int8  `json:"isDelete"`
}

type Documentary struct {
	EpID     int64 `json:"id"`
	SeasonID int64
	AID      int64 `json:"aid"`
	CID      int64 `json:"cid"`
	Status   int8  `json:"is_delete"`
}

type Archive struct {
	AID      int64
	SeasonID int64
	Tp       int8
	Status   int8
	EpID     int64
	CID      int64
}
