package model

import arcmdl "go-common/app/service/main/archive/api"

// Rank bigdata rank struct
type Rank struct {
	Note string         `json:"note"`
	Code int            `json:"code"`
	Page int            `json:"page"`
	Num  int            `json:"num"`
	List []*RankArchive `json:"list"`
}

// RankArchive bigdata rank archive struct
type RankArchive struct {
	Aid         interface{}   `json:"aid"`
	Author      string        `json:"author"`
	Coins       int32         `json:"coins"`
	Duration    string        `json:"duration"`
	Mid         int64         `json:"mid"`
	Pic         string        `json:"pic"`
	Play        interface{}   `json:"play"`
	Pts         int           `json:"pts"`
	Title       string        `json:"title"`
	Trend       *int          `json:"trend"`
	VideoReview int32         `json:"video_review"`
	Rights      arcmdl.Rights `json:"rights"`
	Others      []*Other      `json:"others,omitempty"`
}

// Other bigdata other rank struct
type Other struct {
	Aid         interface{}   `json:"aid"`
	Play        interface{}   `json:"play"`
	VideoReview int32         `json:"video_review"`
	Coins       int32         `json:"coins"`
	Pts         int           `json:"pts"`
	Title       string        `json:"title"`
	Pic         string        `json:"pic"`
	Duration    string        `json:"duration"`
	Rights      arcmdl.Rights `json:"rights"`
}

// RankIndex rank index struct.
type RankIndex struct {
	Code  int                      `json:"code"`
	Pages int                      `json:"pages"`
	Num   int                      `json:"num"`
	List  map[string]*IndexArchive `json:"list"`
}

// IndexArchive rank index archive struct.
type IndexArchive struct {
	Aid         string        `json:"aid"`
	Typename    string        `json:"typename"`
	Title       string        `json:"title"`
	Subtitle    string        `json:"subtitle"`
	Play        interface{}   `json:"play"`
	Review      int32         `json:"review"`
	VideoReview int32         `json:"video_review"`
	Favorites   int32         `json:"favorites"`
	Mid         int64         `json:"mid"`
	Author      string        `json:"author"`
	Description string        `json:"description"`
	Create      string        `json:"create"`
	Pic         string        `json:"pic"`
	Coins       int32         `json:"coins"`
	Duration    string        `json:"duration"`
	Badgepay    bool          `json:"badgepay"`
	Rights      arcmdl.Rights `json:"rights"`
}

// RankRecommend rank recommend data struct
type RankRecommend struct {
	Code  int             `json:"code"`
	Pages int             `json:"pages"`
	Num   int             `json:"num"`
	List  []*IndexArchive `json:"list"`
}

// RankRegion rank region data struct
type RankRegion struct {
	Hot         *RankDetail `json:"hot"`
	HotOriginal *RankDetail `json:"hot_original"`
}

// RankDetail rank region detail struct
type RankDetail struct {
	Note string           `json:"note"`
	Code int              `json:"code"`
	Page int              `json:"page"`
	Num  int              `json:"num"`
	List []*RegionArchive `json:"list"`
}

// RegionArchive bigdata region rank archive struct
type RegionArchive struct {
	Aid         string        `json:"aid"`
	Typename    string        `json:"typename"`
	Title       string        `json:"title"`
	Subtitle    string        `json:"subtitle"`
	Play        interface{}   `json:"play"`
	Review      int32         `json:"review"`
	VideoReview int32         `json:"video_review"`
	Favorites   int32         `json:"favorites"`
	Mid         int64         `json:"mid"`
	Author      string        `json:"author"`
	Description string        `json:"description"`
	Create      string        `json:"create"`
	Pic         string        `json:"pic"`
	Coins       int32         `json:"coins"`
	Duration    string        `json:"duration"`
	Badgepay    bool          `json:"badgepay"`
	Pts         int           `json:"pts"`
	Rights      arcmdl.Rights `json:"rights"`
}

// TagArchive bigdata region rank archive struct
type TagArchive struct {
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	Description string        `json:"description"`
	Pic         string        `json:"pic"`
	Play        string        `json:"play"`
	Favorites   string        `json:"favorites"`
	Mid         string        `json:"mid"`
	Review      string        `json:"review"`
	CreatedAt   string        `json:"created_at"`
	VideoReview string        `json:"video_review"`
	Coins       string        `json:"coins"`
	Duration    string        `json:"duration"`
	Aid         int64         `json:"aid"`
	Pts         int           `json:"pts"`
	Trend       int           `json:"trend"`
	Rights      arcmdl.Rights `json:"rights"`
}

// RankData rank service return data
type RankData struct {
	Note string         `json:"note"`
	List []*RankArchive `json:"list"`
}

// RankNewArchive rank archive new struct
type RankNewArchive struct {
	*NewArchive
	*RankStat
	Others []*NewArchive `json:"others,omitempty"`
}

// RankNew rank new struct.
type RankNew struct {
	Note string            `json:"note"`
	List []*RankNewArchive `json:"list"`
}

// NewArchive new rank archive struct
type NewArchive struct {
	Aid   int64 `json:"aid"`
	Score int   `json:"score"`
}

// RankStat rank archive stat.
type RankStat struct {
	Play  int32 `json:"play"`
	Coin  int32 `json:"coin"`
	Danmu int32 `json:"danmu"`
}

// Custom game custom struct
type Custom struct {
	Aid   int64  `json:"aid"`
	Title string `json:"title"`
	Pic   string `json:"pic"`
	Note  string `json:"note"`
	Pos   int    `json:"-"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}
