package model

// ParamCale  calendar params.
type ParamCale struct {
	Stime int64 `form:"stime" validate:"required"`
	Etime int64 `form:"etime" validate:"required"`
}

// ParamContest matchs params.
type ParamContest struct {
	Mid    int64   `form:"mid" validate:"gte=0"`
	Gid    int64   `form:"gid" validate:"gte=0"`
	Tid    int64   `form:"tid" validate:"gte=0"`
	Stime  string  `form:"stime"`
	Etime  string  `form:"etime"`
	GState string  `form:"g_state"`
	Sids   []int64 `form:"sids,split"`
	Sort   int     `form:"sort"`
	Pn     int     `form:"pn"  validate:"gt=0"`
	Ps     int     `form:"ps"  validate:"gt=0,lte=50"`
}

// ParamVideo video params
type ParamVideo struct {
	Mid  int64 `form:"mid"   validate:"gte=0"`
	Gid  int64 `form:"gid"   validate:"gte=0"`
	Tid  int64 `form:"tid"   validate:"gte=0"`
	Year int64 `form:"year"  validate:"gte=0"`
	Tag  int64 `form:"tag"   validate:"gte=0"`
	Sort int64 `form:"sort"  validate:"gte=0"`
	Pn   int   `form:"pn"    validate:"gt=0"`
	Ps   int   `form:"ps"    validate:"gt=0,lte=50"`
}

// ParamSearch search video params
type ParamSearch struct {
	Pn      int    `form:"pn"    validate:"gt=0"`
	Ps      int    `form:"ps"    validate:"gt=0"`
	Keyword string `form:"keyword" validate:"required"`
	Sort    int64  `form:"sort"  validate:"gte=0"`
}

// ParamSeason season params.
type ParamSeason struct {
	VMID int64 `form:"vmid"`
	Sort int64 `form:"sort"`
	Pn   int   `form:"pn"  validate:"gt=0"`
	Ps   int   `form:"ps"  validate:"gt=0,lte=50"`
}

// ParamFilter  filter video params
type ParamFilter struct {
	Mid   int64  `form:"mid"   validate:"gte=0"`
	Gid   int64  `form:"gid"   validate:"gte=0"`
	Tid   int64  `form:"tid"   validate:"gte=0"`
	Year  int64  `form:"year"  validate:"gte=0"`
	Tag   int64  `form:"tag"   validate:"gte=0"`
	Stime string `form:"stime" `
	Etime string `form:"etime" `
}

// ParamActPoint matchs params.
type ParamActPoint struct {
	Aid  int64 `form:"aid" validate:"gt=0"`
	MdID int64 `form:"md_id" validate:"gt=0"`
	Sort int   `form:"sort"`
	Pn   int   `form:"pn"  validate:"gt=0"`
	Ps   int   `form:"ps"  validate:"gt=0,lte=50"`
}

// ParamActTop matchs params.
type ParamActTop struct {
	Aid   int64  `form:"aid" validate:"gt=0"`
	Sort  int    `form:"sort"`
	Stime string `form:"stime" `
	Etime string `form:"etime" `
	Pn    int    `form:"pn"  validate:"gt=0"`
	Ps    int    `form:"ps"  validate:"gt=0,lte=50"`
}

// ParamFav app fav list.
type ParamFav struct {
	VMID  int64   `form:"vmid"`
	Sids  []int64 `form:"sids,split"`
	Stime string  `form:"stime"`
	Etime string  `form:"etime"`
	Sort  int     `form:"sort"`
	Pn    int     `form:"pn" default:"1" validate:"min=1"`
	Ps    int     `form:"ps" default:"50" validate:"min=1"`
}

// ParamLd leidata param
type ParamLd struct {
	Route string `form:"route"`
}

// ParamCDRecent contest recently match
type ParamCDRecent struct {
	HomeID int64 `form:"home_id" validate:"gt=0"`
	AwayID int64 `form:"away_id" validate:"gt=0"`
	CID    int64 `form:"cid" validate:"gt=0"`
	Ps     int64 `form:"ps" default:"8" validate:"lte=10"`
}

// ParamGame game
type ParamGame struct {
	MatchID int64   `form:"match_id" validate:"required"`
	GameIDs []int64 `form:"game_ids,split" validate:"required"`
	Tp      int64   `form:"tp" default:"1" validate:"min=1"`
}

// ParamLeidas .
type ParamLeidas struct {
	IDs []int64 `form:"ids,split" validate:"required"`
	Tp  int64   `form:"tp" default:"1" validate:"min=1"`
}
