package model

import (
	arcmodel "go-common/app/service/main/archive/model/archive"
	xtime "go-common/library/time"
)

const (
	// TopDataLenth .
	TopDataLenth int = 5
)

// DataType .
type DataType int8

// DataType 数据类型，1累计，2昨日，3上周，4上月
/* ENUM(
Accumulate = 1
Day = 2
Week = 3
Month = 4
ActiveFans = 5
)*/
const (
	// DataTypeAccumulate is a DataType of type Accumulate
	DataTypeAccumulate DataType = 1
	// DataTypeDay is a DataType of type Day
	DataTypeDay DataType = 2
	// DataTypeWeek is a DataType of type Week
	DataTypeWeek DataType = 3
	// DataTypeMonth is a DataType of type Month
	DataTypeMonth DataType = 4
	// DataTypeActiveFans active fans
	DataTypeActiveFans DataType = 5
)

// DataViewTypeSummary .
type DataViewTypeSummary int8

const (
	// SignUpsAccumulate signed up accumulate amount.
	SignUpsAccumulate DataViewTypeSummary = 1
	// FansIncr signed up fans incr amount.
	FansIncr DataViewTypeSummary = 2
	// VideoUpsIncr signed up videoup incr amount.
	VideoUpsIncr DataViewTypeSummary = 3
	// PlaysIncr signed up paly incr amount.
	PlaysIncr DataViewTypeSummary = 4
)

// DataViewFansTop .
type DataViewFansTop int8

const (
	// McnFansIncr .
	McnFansIncr DataViewFansTop = 1
	// McnFansIncrRate .
	McnFansIncrRate DataViewFansTop = 2
	// UpFansIncr .
	UpFansIncr DataViewFansTop = 3
	// UpFansIncrRate .
	UpFansIncrRate DataViewFansTop = 4
)

// MCNDataSummary .
type MCNDataSummary struct {
	ID                       int64      `json:"id"`
	MCNID                    int64      `json:"mcn_mid"`
	SignID                   int64      `json:"sign_id"`
	UPCount                  int64      `json:"up_count"`
	FansCountAccumulate      int64      `json:"fans_count_accumulate"`
	FansCountOnline          int64      `json:"fans_count_online"`
	FansCountReal            int64      `json:"fans_count_real"`
	FansCountCheat           int64      `json:"fans_count_cheat"`
	FansCountCheatAccumulate int64      `json:"fans_count_cheat_accumulate"`
	FansCountIncreaseDay     int64      `json:"fans_count_increase_day"`
	PlayCountAccumulate      int64      `json:"play_count_accumulate"`
	PlayCountIncreaseDay     int64      `json:"play_count_increase_day"`
	ArchiveCountAccumulate   int64      `json:"archive_count_accumulate"`
	ArchiveCountIncreaseDay  int64      `json:"archive_count_increase_day"`
	ActiveTID                int64      `json:"active_tid"`
	GenerateDate             xtime.Time `json:"generate_date"`
	Ctime                    xtime.Time `json:"ctime"`
	Mtime                    xtime.Time `json:"mtime"`
}

// MCNDataUP .
type MCNDataUP struct {
	ID                     int64      `json:"id"`
	MCNID                  int64      `json:"mcn_mid"`
	SignID                 int64      `json:"sign_id"`
	UPMID                  int64      `json:"up_mid"`
	DataType               int8       `json:"data_type"`
	FansCountAll           int64      `json:"fans_count_all"`
	FansCountActive        int64      `json:"fans_count_active"`
	FansIncreaseAccumulate int64      `json:"fans_increase_accumulate"`
	ArchiveCount           int64      `json:"archive_count"`
	PlayCount              int64      `json:"play_count"`
	FansIncreaseMonth      int64      `json:"fans_increase_month"`
	GenerateDate           xtime.Time `json:"generate_date"`
	Ctime                  xtime.Time `json:"ctime"`
	Mtime                  xtime.Time `json:"mtime"`
}

// MCNDataArchiveRank .
type MCNDataArchiveRank struct {
	ID                  int64      `json:"id"`
	MCNID               int64      `json:"mcn_mid"`
	SignID              int64      `json:"sign_id"`
	ArchiveID           int64      `json:"archive_id"`
	ArchiveTitle        string     `json:"archive_title"`
	UPMID               int64      `json:"up_mid"`
	LikeCountAccumulate int64      `json:"like_count_accumulate"`
	LikeCountIncrease   int64      `json:"like_count_increase"`
	PlayCountIncrease   int64      `json:"play_count_increase"`
	DataType            int8       `json:"data_type"`
	TID                 int64      `json:"tid"`
	CtimeArchive        xtime.Time `json:"ctime_archive"`
	Ctime               xtime.Time `json:"ctime"`
	Mtime               xtime.Time `json:"mtime"`
}

// RankDataBase 基本排行信息
type RankDataBase struct {
	Tid      int16    `json:"tid"`
	DataType DataType `json:"data_type"`
}

// TidnameInfo tid name
type TidnameInfo struct {
	Tid  int16  `json:"tid"`
	Name string `json:"name"`
}

// RankArchiveLikeInfo archive like rank info
type RankArchiveLikeInfo struct {
	RankDataBase
	ArchiveID       int64            `json:"archive_id"` // 稿件ID
	ArchiveTitle    string           `json:"archive_title"`
	Pic             string           `json:"pic"` // 封面
	TidName         string           `json:"tid_name"`
	LikesIncrease   int64            `json:"likes_increase"`
	LikesAccumulate int64            `json:"likes_accumulate"`
	PlayIncrease    int64            `json:"play_increase"`
	PlayAccumulate  int64            `json:"play_accumulate"`
	Ctime           xtime.Time       `json:"ctime"`
	Stat            arcmodel.Stat3   `json:"stat"`   // 统计信息
	Author          arcmodel.Author3 `json:"author"` // up主信息
}

// TotalMcnDataInfo .
type TotalMcnDataInfo struct {
	BaseInfo  *McnDataOverview  `json:"base_info"`
	TopInfo   *McnDataTopInfo   `json:"top_info"`
	TypesInfo *McnDataTypesInfo `json:"types_info"`
}

// McnDataTopInfo .
type McnDataTopInfo struct {
	McnFansIncr     []*FansRankIncr  `json:"mcn_fans_incr"`
	McnFansRateIncr []*FansRankIncr  `json:"mcn_fans_rate_incr"`
	UpFansIncr      []*FansRankIncr  `json:"up_fans_incr"`
	UpFansRateIncr  []*FansRankIncr  `json:"up_fans_rate_incr"`
	ArcLikesIncr    []*LikesRankIncr `json:"arc_likes_incr"`
}

// FansRankIncr .
type FansRankIncr struct {
	SignID   int64  `json:"sign_id"`
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	Rank     int16  `json:"rank"`
	FansIncr int64  `json:"fans_incr"`
	Fans     int64  `json:"fans"`
	RateIncr int64  `json:"rate_incr"`
}

// LikesRankIncr .
type LikesRankIncr struct {
	McnMid    int64  `json:"mcn_mid"`
	McnName   string `json:"mcn_name"`
	UpMid     int64  `json:"up_mid"`
	UpName    string `json:"up_name"`
	AVID      int64  `json:"avid"`
	AVTitle   string `json:"av_title"`
	TID       int16  `json:"tid"`
	TypeName  string `json:"type_name"`
	LikesIncr int64  `json:"likes_incr"`
	PlayIncr  int64  `json:"play_incr"`
	SignID    int64  `json:"sign_id"`
}

// McnDataTypesInfo .
type McnDataTypesInfo struct {
	SignUps     []*DataTypes `json:"sign_ups"`
	FansIncr    []*DataTypes `json:"fans_incr"`
	VideoupIncr []*DataTypes `json:"videoup_incr"`
	PlayIncr    []*DataTypes `json:"play_incr"`
}

// DataTypes .
type DataTypes struct {
	TID      int16  `json:"tid"`
	TypeName string `json:"type_name"`
	Total    int64  `json:"total"`
	Amount   int64  `json:"amount"`
	Rate     int64  `json:"rate"`
}

// McnDataOverview base data.
type McnDataOverview struct {
	Mcns        int64 `json:"mcns"`
	SignUps     int64 `json:"sign_ups"`
	SignUpsIncr int64 `json:"sign_ups_incr"`
	Fans50      int64 `json:"fans_50"`
	Fans10      int64 `json:"fans_10"`
	Fans1       int64 `json:"fans_1"`
	FansIncr50  int64 `json:"fans_incr_50"`
	FansIncr10  int64 `json:"fans_incr_10"`
	FansIncr1   int64 `json:"fans_incr_1"`
}

// McnRankFansOverview top5 data.
type McnRankFansOverview struct {
	ID           int64      `json:"id"`
	SignID       int64      `json:"sign_id"`
	Mid          int64      `json:"mid"`
	DataView     int8       `json:"data_view"`
	DataType     int8       `json:"data_type"`
	Rank         int16      `json:"rank"`
	FansIncr     int64      `json:"fans_incr"`
	Fans         int64      `json:"fans"`
	GenerateDate xtime.Time `json:"generate_date"`
	Ctime        xtime.Time `json:"ctime"`
	Mtime        xtime.Time `json:"mtime"`
}

// McnDataTypeSummary tids data.
type McnDataTypeSummary struct {
	ID           int64      `json:"id"`
	Tid          int16      `json:"tid"`
	DataView     int8       `json:"data_view"`
	DataType     int8       `json:"data_type"`
	Amount       int64      `json:"amount"`
	GenerateDate xtime.Time `json:"generate_date"`
	Ctime        xtime.Time `json:"ctime"`
	Mtime        xtime.Time `json:"mtime"`
}

// McnRankArchiveLikesOverview total mcn arc rank likes top
type McnRankArchiveLikesOverview struct {
	ID           int64      `json:"id"`
	McnMid       int64      `json:"mcn_mid"`
	UpMid        int64      `json:"up_mid"`
	SignID       int64      `json:"sign_id"`
	Avid         int64      `json:"avid"`
	Tid          int16      `json:"tid"`
	Rank         int16      `json:"rank"`
	DataType     int8       `json:"data_type"`
	Likes        int64      `json:"likes"`
	Plays        int64      `json:"plays"`
	GenerateDate xtime.Time `json:"generate_date"`
	Ctime        xtime.Time `json:"ctime"`
	Mtime        xtime.Time `json:"mtime"`
}
