package datamodel

import (
	"strconv"
	"time"

	xtime "go-common/library/time"
)

// LogTime
// deserialize format "2006-01-02"
// serialize format unix timestamp
type LogTime time.Time

const (
	timeLayout          = "2006-01-02"
	timeLayoutWithQuote = `"2006-01-02"`
)

//MarshalJSON marshal logdate as timestamp
func (t *LogTime) MarshalJSON() ([]byte, error) {
	var scratch [64]byte
	var dst = strconv.AppendInt(scratch[:0], (*time.Time)(t).Unix(), 10)
	return dst, nil
}

//UnmarshalJSON parse timestamp or something like "2006-01-02"
func (t *LogTime) UnmarshalJSON(b []byte) (err error) {
	var str = string(b)
	if str == "null" {
		return nil
	}

	if len(str) == 0 {
		return
	}

	if str[0] == '"' {
		// Fractional seconds are handled implicitly by Parse.
		tmp, e := time.ParseInLocation(timeLayoutWithQuote, string(b), time.Local)
		if e != nil {
			err = e
			return
		}
		*t = LogTime(tmp)
	} else {
		// parse as timestamp
		num, e := strconv.ParseInt(str, 10, 64)
		if e != nil {
			err = e
			return
		}
		*t = LogTime(time.Unix(num, 0))
	}

	return
}

//Time export library time
func (t *LogTime) Time() xtime.Time {
	return xtime.Time(time.Time(*t).Unix())
}

// see mcn_data.sql for more info

//McnStatisticBaseInfo2 new from data center
type McnStatisticBaseInfo2 struct {
	DanmuAll int64     `json:"danmu_all"`
	DanmuInc int64     `json:"danmu_inc"`
	ReplyAll int64     `json:"reply_all"`
	ReplyInc int64     `json:"reply_inc"`
	ShareAll int64     `json:"share_all"`
	ShareInc int64     `json:"share_inc"`
	CoinAll  int64     `json:"coin_all"`
	CoinInc  int64     `json:"coin_inc"`
	FavAll   int64     `json:"fav_all"`
	FavInc   int64     `json:"fav_inc"`
	LikeAll  int64     `json:"like_all"`
	LikeInc  int64     `json:"like_inc"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"-"`
}

//DmConMcnArchiveD 投稿数及昨日增量
type DmConMcnArchiveD struct {
	McnStatisticBaseInfo2
	ID         int64   `json:"-"`
	SignID     int64   `json:"-"`
	McnMid     int64   `json:"-"`
	LogDate    LogTime `json:"log_date"`
	UpAll      int64   `json:"up_all"`
	ArchiveAll int64   `json:"archive_all"`
	ArchiveInc int64   `json:"archive_inc"`
	PlayAll    int64   `json:"play_all"`
	PlayInc    int64   `json:"play_inc"`
	FansAll    int64   `json:"fans_all"`
	FansInc    int64   `json:"fans_inc"`
}

//DmConMcnIndexIncD 播放/弹幕/评论/分享/硬币/收藏/点赞数每日增量
type DmConMcnIndexIncD struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	Value   int64     `json:"value"`
	Type    string    `json:"-"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnIndexSourceD mcn播放/弹幕/评论/分享/硬币/收藏/点赞来源分区
type DmConMcnIndexSourceD struct {
	ID       int64     `json:"-"`
	SignID   int64     `json:"-"`
	McnMid   int64     `json:"-"`
	LogDate  LogTime   `json:"log_date"`
	TypeID   int64     `json:"type_id"`
	Rank     int64     `json:"rank"`
	Value    int64     `json:"value"`
	Type     string    `json:"-"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"-"`
	TypeName string    `json:"type_name"`
}

//DmConMcnPlaySourceD #mcn稿件播放来源占比
type DmConMcnPlaySourceD struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	Iphone  int64     `json:"iphone"`
	Android int64     `json:"android"`
	Pc      int64     `json:"pc"`
	H5      int64     `json:"h5"`
	Other   int64     `json:"other"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansSexW #游客/粉丝性别占比
type DmConMcnFansSexW struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	Male    int64     `json:"male"`
	Female  int64     `json:"female"`
	Type    string    `json:"-"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansAgeW #游客/粉丝年龄分布
type DmConMcnFansAgeW struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	A       int64     `json:"a"` //'0-16岁人数',
	B       int64     `json:"b"` //'16-25岁人数',
	C       int64     `json:"c"` //'25-40岁人数',
	D       int64     `json:"d"` //'40岁以上人数',
	Type    string    `json:"-"` //'粉丝类型，guest、fans',
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansPlayWayW 游客/粉丝观看途径
type DmConMcnFansPlayWayW struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	App     int64     `json:"app"`
	Pc      int64     `json:"pc"`
	Outside int64     `json:"outside"`
	Other   int64     `json:"other"`
	Type    string    `json:"-"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansAreaW #游客/粉丝地区分布
type DmConMcnFansAreaW struct {
	ID       int64     `json:"-"`
	SignID   int64     `json:"-"`
	McnMid   int64     `json:"-"`
	LogDate  LogTime   `json:"log_date"`
	Province string    `json:"province"`
	User     int64     `json:"user"`
	Type     string    `json:"-"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"-"`
}

//DmConMcnFansTypeW #游客/粉丝倾向分布
type DmConMcnFansTypeW struct {
	ID       int64     `json:"-"`
	SignID   int64     `json:"-"`
	McnMid   int64     `json:"-"`
	LogDate  LogTime   `json:"log_date"`
	TypeID   int64     `json:"type_id"`
	Play     int64     `json:"play"`
	Type     string    `json:"-"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"-"`
	TypeName string    `json:"type_name"`
}

//DmConMcnFansD #mcn粉丝数相关
type DmConMcnFansD struct {
	ID                int64     `json:"-"`
	SignId            int64     `json:"-"`
	McnMid            int64     `json:"-"`
	LogDate           LogTime   `json:"log_date"`
	FansAll           int64     `json:"fans_all"`
	FansInc           int64     `json:"fans_inc"`
	ActFans           int64     `json:"act_fans"`
	FansDecAll        int64     `json:"fans_dec_all"`
	FansDec           int64     `json:"fans_dec"`
	ViewFansRate      float64   `json:"view_fans_rate"`
	ActFansRate       float64   `json:"act_fans_rate"`
	ReplyFansRate     float64   `json:"reply_fans_rate"`
	DanmuFansRate     float64   `json:"danmu_fans_rate"`
	CoinFansRate      float64   `json:"coin_fans_rate"`
	LikeFansRate      float64   `json:"like_fans_rate"`
	FavFansRate       float64   `json:"fav_fans_rate"`
	ShareFansRate     float64   `json:"share_fans_rate"`
	LiveGiftFansRate  float64   `json:"live_gift_fans_rate"`
	LiveDanmuFansRate float64   `json:"live_danmu_fans_rate"`
	Ctime             time.Time `json:"-"`
	Mtime             time.Time `json:"-"`
}

//DmConMcnFansIncD #mcn粉丝按天增量
type DmConMcnFansIncD struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	FansInc int64     `json:"fans_inc"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansDecD #mcn粉丝按天取关数
type DmConMcnFansDecD struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	FansDec int64     `json:"fans_dec"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
}

//DmConMcnFansAttentionWayD #mcn粉丝关注渠道
type DmConMcnFansAttentionWayD struct {
	ID       int64     `json:"-"`
	SignID   int64     `json:"-"`
	McnMid   int64     `json:"-"`
	LogDate  LogTime   `json:"log_date"`
	Homepage int64     `json:"homepage"`
	Video    int64     `json:"video"`
	Article  int64     `json:"article"`
	Music    int64     `json:"music"`
	Other    int64     `json:"other"`
	Ctime    time.Time `json:"-"`
	Mtime    time.Time `json:"-"`
}

// DmConMcnFansTagW #游客/粉丝标签地图分布
type DmConMcnFansTagW struct {
	ID      int64     `json:"-"`
	SignID  int64     `json:"-"`
	McnMid  int64     `json:"-"`
	LogDate LogTime   `json:"log_date"`
	TagID   int64     `json:"tag_id"`
	Play    int64     `json:"play"`
	Type    string    `json:"-"`
	Ctime   time.Time `json:"-"`
	Mtime   time.Time `json:"-"`
	TagName string    `json:"tag_name"`
}
